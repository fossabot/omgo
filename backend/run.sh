#!/bin/bash

cd ./snowflake/
sh ./run.sh
cd ..

cd ./agent/
sh ./run.sh
cd ..

cd ./db/mongodb/
sh ./run.sh
cd ..

cd ./redis/
sh ./run.sh
cd ..

cd ./dbservice/
sh ./run.sh
cd ../../

cd ./reception/
sh ./run.sh
cd ..
