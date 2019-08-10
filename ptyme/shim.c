#define _GNU_SOURCE 600

#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/epoll.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>

#define EPOLL_EVENTS 128
#define PTSNAME_SIZE 1024

static int write_all(int fd, const void *buf, size_t count) {
    const void *p = buf;
    size_t remain = count;
    ssize_t n; 
    while (remain > 0) {
        do {
            n = write(fd, p, remain);
        } while (n == -1 && errno == EINTR);
        
        if (n <= 0) {
            return -1;
        }

        p += n;
        remain -= n;
    }

    return 0;
}

static int socket_listen(char *port) {
    struct addrinfo hints;
    memset(&hints, 0, sizeof (struct addrinfo));
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    hints.ai_flags = AI_PASSIVE;

    struct addrinfo *ainfo;
    if (getaddrinfo(NULL, port, &hints, &ainfo) != 0) {
        perror("getaddrinfo() failed");
        return -1;
    }

    int sock;
    struct addrinfo *rp;
    for (rp = ainfo; rp != NULL; rp = rp->ai_next) {
        sock = socket(rp->ai_family, rp->ai_socktype, rp->ai_protocol);
        if (sock == -1) {
            continue;
        }

        if (bind(sock, rp->ai_addr, rp->ai_addrlen) == 0) {
            break;
        }

        perror("bind() failed");
        close(sock);
    }
    freeaddrinfo(ainfo);

    if (rp == NULL) {
        return -1;
    }

    if (listen(sock, SOMAXCONN) != 0) {
        return -1;
    }

    return sock;
}

struct pt_info {
    int master_fd;
    char slave_name[PTSNAME_SIZE];
};

static int epoll_add_fd(int epoll, int fd) {
    struct epoll_event ev;
    ev.data.fd = fd;
    ev.events = EPOLLIN;
    return epoll_ctl(epoll, EPOLL_CTL_ADD, fd, &ev);
}

// atsock - keeps track of the attached sockets
struct atsock {
    int fd;
    struct atsock *next;
};

struct atsock *atsock_new(int conn_fd) {
    struct atsock *s = malloc(sizeof(struct atsock));
    if (s != NULL) {
        s->fd = conn_fd;
        s->next = NULL;
    }
    return s;
}

struct atsock *atsock_save(struct atsock *head, struct atsock *new) {
    assert(new != NULL);
    if (head == NULL) {
        return new;
    }

    struct atsock *cur = head;
    while (cur->next != NULL) {
        cur = cur->next;
    }
    cur->next = new;
    return head;
}

struct atsock *atsock_erase(struct atsock *head, int conn_fd) {
    assert(head != NULL);
    if (head->fd == conn_fd) {
        struct atsock *next = head->next;
        free(head);
        return next;
    }

    struct atsock *prev = head, *cur = head->next;
    while (cur) {
        if (cur->fd != conn_fd) {
            prev = cur;
            cur = cur->next;
            continue;
        }

        prev->next = cur->next;
        free(cur);
        return head;
    }

    assert(0);  // s not found in the list
    return NULL;
}

static void run_master(struct pt_info *pti) {
    assert(pti != NULL);

    int exit_code = EXIT_FAILURE;
    int attach_sock = -1;
    int epoll = epoll_create1(0);
    if (epoll == -1) {
        perror("epoll_create1() failed");
        goto exit;
    }

    if (epoll_add_fd(epoll, pti->master_fd) != 0) {
        perror("epoll_add_fd(master_fd) failed");
        goto exit;
    }

    attach_sock = socket_listen("43210");
    if (attach_sock == -1) {
        perror("socket_listen() failed");
        goto exit;
    }

    if (epoll_add_fd(epoll, attach_sock) != 0) {
        perror("epoll_add_fd(attach_sock) failed");
        goto exit;
    }
    
    struct atsock *head = NULL;
    struct epoll_event evlist[EPOLL_EVENTS];
    while (1) {
        int nready = epoll_wait(epoll, evlist, EPOLL_EVENTS, -1);
        if (nready == -1 && errno == EINTR) {
            continue;
        }
        if (nready == -1) {
            perror("epoll_wait() failed");
            break;
        }

        for (int i = 0; i < nready; i++) {
            int fd = evlist[i].data.fd;
            if (evlist[i].events & EPOLLIN) {
                if (fd == pti->master_fd) {
                    // read from pty and forward data to each attached socket
                    char buf[1024];
                    int nread = read(fd, buf, 1023);
                    struct atsock *cur = head;
                    while (nread && cur) {
                        write_all(cur->fd, buf, nread);
                        cur = cur->next;
                    }
                } else if (fd == attach_sock) {
                    int conn;
                    do {
                        conn = accept(fd, NULL, NULL);
                    } while (conn == -1 && errno == EINTR);

                    if (conn == -1) {
                        perror("accept() failed");
                        goto exit;
                    }

                    head = atsock_save(head, atsock_new(conn));
                    if (epoll_add_fd(epoll, conn) != 0) {
                        perror("epoll_add_fd(conn) failed");
                        goto exit;
                    }
                    printf("accepted new sock conn\n");
                } else {
                    // read from attached socket and forward to pty
                    char buf[1024];
                    int nread = read(fd, buf, 1023);
                    if (nread == 0) {
                        head = atsock_erase(head, fd);
                        if (epoll_ctl(epoll, EPOLL_CTL_DEL, fd, NULL) != 0) {
                            perror("epoll_ctl(EPOLL_CTL_DEL) failed");
                            goto exit;
                        }
                        printf("disconnected sock\n");
                    } else {
                        if (write_all(pti->master_fd, buf, nread) != 0) {
                            perror("write_all(master_fd) failed");
                            goto exit;
                        }
                    }
                }
            } else if (evlist[i].events & (EPOLLHUP | EPOLLERR)) {
                if (fd == attach_sock) {
                    attach_sock = -1;
                    if (epoll_ctl(epoll, EPOLL_CTL_DEL, fd, NULL) != 0) {
                        perror("epoll_ctl(EPOLL_CTL_DEL, attach_sock) failed");
                        goto exit;
                    }
                    printf("attach_sock failed\n");
                } else if (fd == pti->master_fd) {
                    exit_code = EXIT_SUCCESS;
                    goto exit;
                } else {
                    head = atsock_erase(head, fd);
                    if (epoll_ctl(epoll, EPOLL_CTL_DEL, fd, NULL) != 0) {
                        perror("epoll_ctl(EPOLL_CTL_DEL) failed");
                        goto exit;
                    }
                    printf("disconnected sock\n");
                }
            }
        }
    }

    exit_code = EXIT_SUCCESS;

exit:
    close(pti->master_fd);
    while (head) {
        close(head->fd);
        head = atsock_erase(head, head->fd);
    }
    if (epoll != -1) {
        close(epoll);
    }
    if (attach_sock != -1) {
        close(attach_sock);
    }
    exit(exit_code);
}

static void run_slave(struct pt_info *pti) {
    assert(pti != NULL);
    close(pti->master_fd);
    setsid();

    int fds = open(pti->slave_name, O_RDWR);
    if (fds >= 0) {
        dup2(fds, 0);
        dup2(fds, 1);
        dup2(fds, 2);
        close(fds);
        execl("/bin/bash", "/bin/bash", NULL);
    } else {
        perror("open(pts_name) failed");
    }
    _exit(127);
}

static int create_pt(struct pt_info *p) {
    errno = 0;

    do {
        if (p == NULL) {
            errno = EINVAL;
            break;
        }

        p->master_fd = posix_openpt(O_RDWR);
        if (p->master_fd < 0) {
            perror("posix_openpt() failed");
            break;
        }
        if (grantpt(p->master_fd) != 0) {
            perror("grantpt() failed");
            break;
        }
        if (unlockpt(p->master_fd) != 0) {
            perror("unlockpt() failed");
            break;
        }

        if (ptsname_r(p->master_fd, p->slave_name, PTSNAME_SIZE) != 0) {
            perror("ptsname_r() failed");
            break;
        }
    } while(0);

    if (errno && p && p->master_fd >= 0) {
        close(p->master_fd);
    }

    return errno;
}

int main() {
    int pid = fork();
    if (pid < 0) {
        perror("fork() failed");
        exit(EXIT_FAILURE);
    } else if (pid) {
        exit(EXIT_SUCCESS);
    }

    int devnull = open("/dev/null", O_RDWR | O_CLOEXEC);
    if (devnull < 0) {
        perror("open('dev/null') failed");
        exit(EXIT_FAILURE);
    }
    if (dup2(devnull, STDIN_FILENO) < 0) {
        perror("dup2(devnul, STDIN) failed");
        exit(EXIT_FAILURE);
    }
    if (dup2(devnull, STDOUT_FILENO) < 0) {
        perror("dup2(devnul, STDOUT) failed");
        exit(EXIT_FAILURE);
    }
    if (dup2(devnull, STDERR_FILENO) < 0) {
        perror("dup2(devnul, STDERR) failed");
        exit(EXIT_FAILURE);
    }

    setsid();

    struct pt_info pti;
    if (create_pt(&pti) != 0) {
        exit(EXIT_FAILURE);
    }

    pid = fork();
    if (pid < 0) {
        perror("fork() failed");
        exit(EXIT_FAILURE);
    } else if (pid == 0) {
        run_slave(&pti);
    } else {
        run_master(&pti);
    }
    return 0;
}

