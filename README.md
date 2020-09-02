你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# conman - [the] container manager

The aim of the project is to implement yet another _container manager_. Primarily, for the sake of [self-]education.
The _conman_ project is heavily inspired by <a href="https://github.com/cri-o/cri-o">cri-o</a> and the ultimate goal is to
make _conman_ <a href="https://github.com/kubernetes/cri-api/">CRI</a>-compatible. This will allow to deploy a Kubernetes
cluster with _conman_ as a container runtime server.

Read more about the project in <a href="https://iximiuz.com/en/posts/conman-the-container-manager-inception/">this article</a>.

## State of the project
Under active development. Not even close to _0.1.0_.

## Run it
So far the only tested platform is CentOS 7 with `go version go1.13.6 linux/amd64`.

```bash
git clone https://github.com/iximiuz/conman.git
cd conman

# Build daemon and client
make bin/conmand
make bin/conmanctl

# Run daemon
sudo bin/conmand

# Prepare dev data
make test/data/rootfs_alpine

# Create containers
sudo bin/conmanctl container create --image test/data/rootfs_alpine/ cont1 -- sleep 100
sudo bin/conmanctl container create --image test/data/rootfs_alpine/ cont2 -- sleep 200

# List containers
sudo bin/conmanctl container list

# Start container 
sudo bin/conmanctl container start <container_id>

# Stop container 
sudo bin/conmanctl container stop <container_id>

# Request container status
sudo bin/conmanctl container status <container_id>

# Remove container 
sudo bin/conmanctl container remove <container_id>
```

## Test it
```bash
# Unit (not really) tests
sudo PATH=/usr/local/bin:$PATH make testunit

# Functional tests
# install jq `yum install jq`
# install bats https://github.com/bats-core/bats-core 
sudo PATH=/usr/local/bin:$PATH make testfunctional

# OCI runtime shim integration tests
make testshimmy
```

## TODO:
- acceptance tests
- shim
  - interactive containers (exec, stdin/stdout support)
  - PTY-controlled containers (eg. shell)
  - attach to a running container

