# Extend vert.x image
FROM java:8-jre

ENV VERTICLE_FILE omgo-webservice-1.0-SNAPSHOT-fat.jar
ENV CONFIG_FILE config.json

# Set the location of the verticles
ENV VERTICLE_HOME /usr/verticles

EXPOSE 60001

# Copy your verticle to the container
COPY target/$VERTICLE_FILE $VERTICLE_HOME/
COPY $CONFIG_FILE $VERTICLE_HOME/

# Launch the verticle
WORKDIR $VERTICLE_HOME
ENTRYPOINT ["sh", "-c"]
CMD ["exec java -jar -Dvertx.metrics.options.enabled=true $VERTICLE_FILE -conf $CONFIG_FILE"]
