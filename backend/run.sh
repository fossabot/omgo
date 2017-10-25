#!/bin/bash

###########################################################
# 3rd party
###########################################################

cd ./3rdparty

# DB mysql
cd ./mongodb
sh ./run.sh $1
cd ..

# DB redis
cd ./redis
sh ./run.sh $1
cd ..

# ETCD
cd ./etcd
sh ./run.sh $1
cd ..

# nsq
cd ./nsq
sh ./nsqlookupd.sh $1
sh ./nsqd.sh $1
cd ..

cd ..

###########################################################
# services
###########################################################

# UUID generator, and ETCD
cd ./snowflake/
sh ./run.sh $1
cd ..

# Data service
cd ./dataservice/
sh ./run.sh $1
cd ..

# Web service
cd ./webservice/
sh ./run.sh $1
cd ..
