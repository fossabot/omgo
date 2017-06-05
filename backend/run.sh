#!/bin/bash

echo starting snowflake
cd ./snowflake/
govendor update +v
sh ./run.sh
cd ..

echo starting agent
cd ./agent/
govendor update +v
sh ./run.sh
cd ..

echo starting mongodb
cd ./db/mongodb/
sh ./run.sh
cd ..

echo starting redis
cd ./redis/
sh ./run.sh
cd ..

echo starting dbservice
cd ./dbservice/
govendor update +v
sh ./run.sh
cd ../../

echo starting reception
cd ./reception/
govendor update +v
sh ./run.sh
cd ..
