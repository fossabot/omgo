#!/bin/bash

# UUID generator, and ETCD
cd ./snowflake/
sh ./run.sh
cd ..

# DB mysql
cd ./db/mysql/
sh ./run.sh
cd ..

# DB redis
cd ./redis/
sh ./run.sh
cd ..

exit

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
