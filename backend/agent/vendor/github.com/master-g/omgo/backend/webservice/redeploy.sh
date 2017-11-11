#!/usr/bin/env bash

mvn package dependency:copy-dependencies
java -jar -Dvertx.metrics.options.enabled=true \
    target/omgo-webservice-1.0-SNAPSHOT-fat.jar \
    -conf src/main/resources/config.json

#export LAUNCHER="io.vertx.core.Launcher"
#export VERTICLE="com.omgo.webservice.MainVerticle"
#export CMD="mvn compile"
#export VERTX_CMD="run"

#mvn compile dependency:copy-dependencies
#java \
#  -cp  $(echo target/dependency/*.jar | tr ' ' ':'):"target/classes" \
#  ${LAUNCHER} ${VERTX_CMD} ${VERTICLE} \
#  --redeploy="src/main/**/*" --on-redeploy="${CMD}" \
#  --launcher-class=${LAUNCHER} \
#  -conf src/main/resources/config.json \
#  -Dvertx.metrics.options.enabled=true \
#  $@
