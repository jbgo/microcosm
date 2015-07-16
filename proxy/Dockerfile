FROM golang:1.4.2-onbuild
ENV DOCKER_HOST=unix:///var/run/docker.sock
VOLUME /var/run/docker.sock:/var/run/docker.sock
CMD ["go-wrapper", "run", "app"]
