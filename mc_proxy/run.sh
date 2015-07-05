#!/bin/bash -e
docker build -t mc_proxy .
docker run --rm -it -e DOCKER_HOST=unix:///var/run/docker.sock -v /var/run/docker.sock:/var/run/docker.sock --volumes-from=mc_haproxy mc_proxy
