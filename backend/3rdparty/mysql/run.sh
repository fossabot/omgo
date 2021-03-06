#!/bin/bash

SID=mysqldb-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

if [ "$1" = "rebuild" ]
then
    docker pull mysql
fi

docker rm -f ${SID}
docker run -d --name ${SID} ${NETHOST} \
  -p 3306:3306 \
  -v ${PWD}/db:/var/lib/mysql \
  -v ${PWD}/docker-entrypoint-initdb.d:/docker-entrypoint-initdb.d \
  -e MYSQL_USER=mysql \
  -e MYSQL_PASSWORD=mysql \
  -e MYSQL_DATABASE=sample \
  -e MYSQL_ROOT_PASSWORD=supersecret \
  mysql
