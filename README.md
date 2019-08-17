# conman - [the] container manager

The aim of the project is to implement yet another _container manager_. Primarily, for the sake of [self-]education.
The _conman_ project is heavily inspired by <a href="https://github.com/cri-o/cri-o">cri-o</a> and the ultimate goal is to
make _conman_ <a href="https://github.com/kubernetes/cri-api/">CRI</a>-compatible. This will allow to deploy a Kubernetes
cluster with _conman_ as a container runtime server.

## State of the project
Under active development. Not even close to _0.1.0_.

## Run it
So far the only tested platform is CentOS 7.

```
git clone https://github.com/iximiuz/conman.git
cd conman

# Build daemon and client
make bin/conmand
make bin/conmanctl

# Run daemon
sudo bin/conmand

# Create a container
mkdir alpine
sudo bash -c "docker export $(docker create alpine) | tar -C alpine/ -xvf -"
sudo bin/conmanctl container create --image `pwd`/alpine mycont
```
