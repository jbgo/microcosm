FROM microcosm/base
WORKDIR /go/src/github.com/jbgo/microcosm/agent
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["sh", "-c", "/entrypoint.sh"]
