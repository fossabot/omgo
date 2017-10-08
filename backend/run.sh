#!/bin/bash

# UUID generator, and ETCD
cd ./snowflake/
sh ./run.sh
cd ..

# DB mysql
cd ./db/mongodb
sh ./run.sh
cd ..

# DB redis
cd ./redis/
sh ./run.sh
cd ../../

# Data service
cd ./dataservice/
sh ./run.sh
cd ..

exit

# DB service
cd ./game/
sh ./run.sh
cd ..

# Agent
cd ./agent/
sh ./run.sh
cd ..

# Reception
cd ./reception/
sh ./run.sh
cd ..
