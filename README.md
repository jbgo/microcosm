# Microcosm

A microservices-aware, production-like development environment for your organization.

### Coming Fall 2015

This project is still in the exploratory phase. Nothing is stable.

### Development setup

Microcosm is built on [go 1.4](http://golang.org/doc/install) and [docker 1.7](https://docs.docker.com/userguide/).
If you're hacking on this project, you will also want to install [godep](https://github.com/tools/godep).
For your convenience, a Vagrantfile is provided.

```
git clone git@github.com:jbgo/microcosm.git
cd microcosm
vagrant up
```

After a few minutes, you will have a vagrant VM running with the IP address 192.168.33.33.
To run docker commands from the host machine, run the following command on your host to set the proper docker configuration environment variables.

```
source env.sh
```

You should now be able to run `docker info` on either the host or the guest and get a successful result.
