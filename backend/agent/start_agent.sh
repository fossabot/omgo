#!/usr/bin/env bash

# --net=host, bind the container to the host interface
# so other container can access this container via localhost:port
# ps. --net=host is not working on macOS, see https://github.com/docker/for-mac/issues/68

HOST=127.0.0.1
SID=agent1

docker rm -f ${SID}
docker build --no-cache --rm=true -t agent .
docker run --name ${SID} -e SERVICE_ID=${SID} --net=host -p 8888:8888 -d -P agent \
    -l :8888 \
    -e http://${HOST}:2379

# register service
curl -L -X PUT http://${HOST}:2379/v2/keys/backends/agent/${SID} -d value=${HOST}:8888
