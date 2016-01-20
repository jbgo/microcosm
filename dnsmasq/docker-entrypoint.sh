#!/bin/bash

echo "address=/docker/$DOCKER_IP" > /etc/dnsmasq.d/docker.conf

exec "$@"
