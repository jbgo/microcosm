FROM microcosm/base
WORKDIR /go/src/github.com/jbgo/microcosm/proxy
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["sh", "-c", "/entrypoint.sh"]
