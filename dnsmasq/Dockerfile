FROM debian:jessie

RUN apt-get update \
  && apt-get install -yq dnsmasq \
  && chmod +s /usr/sbin/dnsmasq \
  && chmod +s /sbin/start-stop-daemon \
  && mkdir -p /etc/dnsmasq.d

COPY ./dnsmasq.conf /etc/dnsmasq.conf

COPY ./docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

ENV DOCKER_IP ${DOCKER_IP}

EXPOSE 53

CMD ["dnsmasq", "-d"]
