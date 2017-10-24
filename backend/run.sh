#!/bin/bash

# UUID generator, and ETCD
cd ./snowflake/
sh ./run.sh $1
cd ..

# DB mysql
cd ./db/mongodb
sh ./run.sh $1
cd ..

# DB redis
cd ./redis/
sh ./run.sh $1
cd ../../

# Data service
cd ./dataservice/
sh ./run.sh $1
cd ..

# Web service
cd ./webservice/
sh ./run.sh $1
cd ..
