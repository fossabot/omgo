#!/bin/bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
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
docker build --no-cache --rm=true -t snowflake .
docker run --rm -d ${NETHOST} -p 40001:40001 \
    --name ${SNOWFLAKE_SID} \
    -e SERVICE_ID=${SNOWFLAKE_SID} \
    -e MACHINE_ID=1 \
    --entrypoint /go/bin/snowflake \
    snowflake \
    -p 40001 \
    -e http://${IPADDR}:2379

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/snowflake/${SNOWFLAKE_SID} -d value=${IPADDR}:40001

# init etcd variables
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/test_key -d value="0"
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/userid -d value="0"
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/snowflake-uuid -d value="0"
