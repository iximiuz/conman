# conman - [the] container manager

The aim of the project is to implement yet another _container manager_. Primarily, for the sake of [self-]education.
The _conman_ project is heavily inspired by <a href="https://github.com/cri-o/cri-o">cri-o</a> and the ultimate goal is to
make _conman_ <a href="https://github.com/kubernetes/cri-api/">CRI</a>-compatible. This will allow to deploy a Kubernetes
cluster with _conman_ as a container runtime server.

## State of the project
Under active development. Not even close to _0.1.0_.

## Run it
So far the only tested platform is CentOS 7.

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
- pick up state of the containers on conman [re]start
- acceptance tests
- shim
  - track container exit code
  - interactive containers (exec, stdin/stdout support)
  - PTY-controlled containers (eg. shell)
  - attach to a running container

