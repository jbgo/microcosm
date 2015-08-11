# Microcosm

A microservices-aware, production-like development environment for your organization.

### Coming Fall 2015

This project is still in the exploratory phase. Nothing is stable.

### Roadmap

* service-centric, because that's what you're building. containers are only an implementation detail
* convention over configuration so you can stop tweaking your environment and just write code
* automatic discovery, classification, and load-balancing of services launched with docker compose
* shared database setups and migrations to save memory on your laptop
* plugin-architecture to run custom tasks based on docker or filesystem events
* launch new services directly from a GitHub repo with a docker compose file
* dashboard to monitor and manage your services and the docker host

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
