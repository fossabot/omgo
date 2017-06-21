#!/bin/bash

# UUID generator, and ETCD
cd ./snowflake/
sh ./run.sh
cd ..

# DB mongodb
cd ./db/mongodb/
sh ./run.sh
cd ..

# DB redis
cd ./redis/
sh ./run.sh
cd ..

# DB service
cd ./dbservice/
sh ./run.sh
cd ../../

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
