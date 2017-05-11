#!/bin/bash

SID=mongodb1

docker rm -f ${SID}

docker run -d --name ${SID} \
  -p 37017:27017 \
  -v ${PWD}/mongod.conf:/etc/mongod.conf \
  -v ${PWD}/db:/data/db \
  mongo \
  --config /etc/mongod.conf
