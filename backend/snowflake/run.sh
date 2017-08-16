#!/bin/bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
SERVICE_NAME=snowflake
SERVICE_PORT=40001
SNOWFLAKE_SID=snowflake-0
ETCD_SID=etcd0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

# ETCD
docker rm -f etcd
docker run --rm -d ${NETHOST} -p 2379:2379 -p 2380:2380 \
    --name etcd \
    quay.io/coreos/etcd \
    /usr/local/bin/etcd \
    --name ${ETCD_SID} \
    --data-dir=data.etcd \
    --initial-advertise-peer-urls http://${LOCALHOST}:2380 \
    --listen-peer-urls http://${LOCALHOST}:2380 \
    --advertise-client-urls http://0.0.0.0:2379 \
    --listen-client-urls http://0.0.0.0:2379 \
    --initial-cluster ${ETCD_SID}=http://${LOCALHOST}:2380 \
    --initial-cluster-state new \
    --initial-cluster-token my-etcd-token

# Snowflake
docker rm -f ${SNOWFLAKE_SID}
docker build --no-cache --rm=true -t ${SERVICE_NAME} .
docker run --rm -d ${NETHOST} -p ${SERVICE_PORT}:${SERVICE_PORT} \
    --name ${SNOWFLAKE_SID} \
    -e SERVICE_ID=${SNOWFLAKE_SID} \
    -e MACHINE_ID=1 \
    --entrypoint /go/bin/${SERVICE_NAME} \
    ${SERVICE_NAME} \
    --service-key backends/${SERVICE_NAME}/${SNOWFLAKE_SID} \
    --service-host ${IPADDR}:${SERVICE_PORT} \
    -p ${SERVICE_PORT} \
    -e http://${IPADDR}:2379


# Deprecated for ETCD v3

# register service
#curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/${SERVICE_NAME}/${SNOWFLAKE_SID} -d value=${IPADDR}:${SERVICE_PORT}

# init etcd variables
#curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/test_key -d value="0"
#curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/snowflake-uuid -d value="0"

# DANGER !!! THIS WILL RESET USER ID !!!
#curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/userid -d value="10000"
