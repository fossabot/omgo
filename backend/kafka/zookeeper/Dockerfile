FROM ubuntu:trusty
RUN apt-get update && apt-get dist-upgrade -y
RUN apt-get install -y zookeeper
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN sed -i 's/ROLLINGFILE/CONSOLE/' /etc/zookeeper/conf/environment
COPY docker-entrypoint.sh /usr/local/bin/
CMD ["/usr/local/bin/docker-entrypoint.sh"]
