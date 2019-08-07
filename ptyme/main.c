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
// #include <termios.h>
#include <unistd.h>

#define MAX_EVENTS 128
#define PTSNAME_SIZE 1024

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
    int epoll = -1;
    int attach_sock = -1;

    epoll = epoll_create1(0);
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
    struct epoll_event evlist[MAX_EVENTS];
    while (1) {
        int nready = epoll_wait(epoll, evlist, MAX_EVENTS, -1);
        printf("nready=%d\n", nready);
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
                    char buf[1024];
                    int nread = read(fd, buf, 1023);
                    for (int j = 0; j < nread; j++) {
                        printf("[%x] %c", buf[j], buf[j]);
                    }
                    printf("\n");
                    // read from pty and forward data to each attached socket
                } else if (fd == attach_sock) {
                    int conn = accept(fd, NULL, NULL);
                    if (conn == -1) {
                        if (errno == EINTR) {
                            continue;
                        }
                        perror("accept() failed");
                    } else {
                        head = atsock_save(head, atsock_new(conn));
                        if (epoll_add_fd(epoll, conn) != 0) {
                            perror("epoll_add_fd(conn) failed");
                            goto exit;
                        }
                        printf("accepted new sock conn\n");
                    }
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
                    }
                }
            }

            if (evlist[i].events & (EPOLLHUP | EPOLLERR)) {
                if (fd != pti->master_fd && fd != attach_sock) {
                    head = atsock_erase(head, fd);
                    if (epoll_ctl(epoll, EPOLL_CTL_DEL, fd, NULL) != 0) {
                        perror("epoll_ctl(EPOLL_CTL_DEL) failed");
                        goto exit;
                    }
                    printf("disconnected sock\n");
                }
                // TODO: handle error
            }
        }
    }

    exit_code = EXIT_SUCCESS;

exit:
    close(pti->master_fd);
    if (epoll != -1) {
        close(epoll);
    }
    if (attach_sock != -1) {
        close(attach_sock);
    }
    exit(exit_code);
}

    // struct termios raw;
    // tcgetattr(STDIN_FILENO, &raw);
    // raw.c_lflag &= ~(ECHO);
    // tcsetattr(STDIN_FILENO, TCSAFLUSH, &raw);

    // char buf[1024];
    // while (1) {
    //     int n = read(0, buf, 1023);
    //     if (n < 0) {
    //         perror("read(STDIN) failed");
    //         break;
    //     }

    //     if (n) {
    //         int nw = write(pti->master_fd, buf, n);
    //         if (nw < 0) {
    //             perror("write(fdm) failed");
    //             break;
    //         }

    //         int nr = read(pti->master_fd, buf, 1023);
    //         if (nr < 0) {
    //             perror("read(fdm) failed");
    //             break;
    //         }
    //         if (nr) {
    //             int nw = write(1, buf, nr);
    //             if (nw < 0) {
    //                 perror("write(STDOUT) failed");
    //                 break;
    //             }
    //         }
    //     }
    // }

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
    printf("ptyme 0.1.0!\n");

    struct pt_info pti;
    if (create_pt(&pti) != 0) {
        exit(EXIT_FAILURE);
    }

    int pid = fork();
    if (pid < 0) {
        perror("fork() failed");
        exit(EXIT_FAILURE);
    }

    if (pid == 0) {
        run_slave(&pti);
    } 

    run_master(&pti);
    return 0;
}

