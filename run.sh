#!/bin/bash

set -e

echo '===== build image: microcosm/base ====='
docker build -t microcosm/base:latest .

echo '===== build image: microcosm/agent ====='
docker build -t microcosm/agent:latest ./agent

echo '===== build image: microcosm/proxy ====='
docker build -t microcosm/proxy:latest ./proxy

echo '===== create container: microcosm-code ====='
docker create --name=microcosm-code \
  -v /home/vagrant/go/src/github.com/jbgo/microcosm:/go/src/github.com/jbgo/microcosm \
  microcosm/base

echo '===== create container: microcosm-configure-proxy ====='
docker create --name=microcosm-configure-proxy \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --volumes-from=microcosm-code \
  --label=microcosm.service=microcosm \
  --label=microcosm.type=task \
  microcosm/proxy

echo '===== create and run container: microcosm-agent ====='
docker run -d --name=microcosm-agent \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --volumes-from=microcosm-code \
  --label=microcosm.service=microcosm \
  --label=microcosm.type=daemon \
  microcosm/agent
