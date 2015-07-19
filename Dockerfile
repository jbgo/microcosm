FROM golang:1.4.2

RUN mkdir -p /go/src
COPY . /go/src/github.com/jbgo/microcosm
WORKDIR /go/src/github.com/jbgo/microcosm

RUN go get github.com/tools/godep
RUN godep restore

ENV DOCKER_HOST=unix:///var/run/docker.sock
ENV GOLANG_VERSION 1.4.2
ENV GOPATH=/go
ENV PATH /usr/src/go/bin:$PATH
ENV PATH /go/bin:$PATH

VOLUME .:/go/src/github.com/jbgo/microcosm

CMD ["echo", "microcosm/base"]
