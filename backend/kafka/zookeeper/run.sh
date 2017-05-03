#!/bin/bash

IPADDR=127.0.0.1
NETHOST=--net=host
PORTMAP1=(2181:2181 2182:2181 2183:2181)
PORTMAP2=(2888:2888 2889:2889 2890:2890)
PORTMAP3=(3888:3888 3889:3889 3890:3890)

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker build -t endocode/zookeeper .

for i in {1..3}
do
    echo ${PORTMAP1[$i - 1]}
    echo ${PORTMAP2[$i - 1]}
    echo ${PORTMAP3[$i - 1]}

    docker rm -f zookeeper${i}
    docker run --rm -d ${NETHOST} \
        --name zookeeper${i} \
        -p ${PORTMAP1[$i - 1]} \
        -p ${PORTMAP2[$i - 1]} \
        -p ${PORTMAP3[$i - 1]} \
        -e ZK_SERVERS="server.1=${IPADDR}:2888:3888 server.2=${IPADDR}:2889:3889 server.3=${IPADDR}:2890:3890" \
        -e ZK_ID=${i} \
        endocode/zookeeper
done
