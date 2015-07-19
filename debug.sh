docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --volumes-from=microcosm-code \
  --volumes-from=microcosm-proxy-data \
  microcosm/proxy
