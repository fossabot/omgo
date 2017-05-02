#!/usr/bin/env bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
SID=snowflake1
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker rm -f etcd
docker run -d ${NETHOST} -p 2379:2379 -p 2380:2380 --name etcd quay.io/coreos/etcd \
    /usr/local/bin/etcd \
    --name etcd0 \
    --data-dir=data.etcd \
    --initial-advertise-peer-urls http://${LOCALHOST}:2380 \
    --listen-peer-urls http://${LOCALHOST}:2380 \
    --advertise-client-urls http://0.0.0.0:2379 \
    --listen-client-urls http://0.0.0.0:2379 \
    --initial-cluster etcd0=http://${LOCALHOST}:2380 \
    --initial-cluster-state new \
    --initial-cluster-token my-etcd-token

docker rm -f ${SID}
docker build --no-cache --rm=true -t snowflake .
docker run --name ${SID} -e SERVICE_ID=${SID} -e MACHINE_ID=1 ${NETHOST} -p 40001:40001 -d -P snowflake \
    -p 40001 \
    -e http://${IPADDR}:2379

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/snowflake/${SID} -d value=${IPADDR}:40001

# init etcd variables
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/test_key -d value="0"
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/userid -d value="0"
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/seqs/snowflake-uuid -d value="0"
