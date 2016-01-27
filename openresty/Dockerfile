FROM debian:jessie
MAINTAINER Jordan Bach <jordanbach@gmail.com>

WORKDIR /tmp

RUN apt-get update
RUN apt-get install -yq libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make build-essential curl
RUN curl -O https://openresty.org/download/ngx_openresty-1.9.7.2.tar.gz
RUN tar xvf ngx_openresty-1.9.7.2.tar.gz

WORKDIR /tmp/ngx_openresty-1.9.7.2

RUN ./configure --with-pcre-jit --with-ipv6
RUN make
RUN make install

RUN mkdir -p /microcosm/logs /microcosm/conf

WORKDIR /microcosm

COPY nginx.conf /microcosm/conf/nginx.conf

VOLUME ["/microcosm"]

EXPOSE 80

CMD ["/usr/local/openresty/nginx/sbin/nginx", "-p", "/microcosm", "-c", "conf/nginx.conf"]
