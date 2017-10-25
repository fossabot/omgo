#!/bin/bash

###########################################################
# 3rd party
###########################################################

# DB mysql
cd ./3rdparty/mongodb
sh ./run.sh $1
cd ../../

# DB redis
cd ./3rdparty/redis/
sh ./run.sh $1
cd ../../

# ETCD
cd ./3rdparty/etcd/
sh ./run.sh $1
cd ../../

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
