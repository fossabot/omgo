#!/usr/bin/env bash

# --net=host, bind the container to the host interface
# so other container can access this container via localhost:port
# ps. --net=host is not working on macOS, see https://github.com/docker/for-mac/issues/68

IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
LOCALHOST=127.0.0.1
SID=agent1

docker rm -f ${SID}
docker build --no-cache --rm=true -t agent .
docker run --name ${SID} -e SERVICE_ID=${SID} -e MACHINE_ID=1 -p 8888:8888 -d -P agent \
    -l :8888 \
    -e http://${IPADDR}:2379

# register service
curl -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/${SID} -d value=${LOCALHOST}:8888
