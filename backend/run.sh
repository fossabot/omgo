#!/bin/bash

cd ./snowflake/
govendor update +v
exec ./run.sh
cd ..

cd ./agent/
govendor update +v
exec ./run.sh
cd ..

cd ./db/mongodb/
exec ./run.sh
cd ..

cd ./redis/
exec ./run.sh
cd ..

cd ./dbservice/
govendor update +v
exec ./run.sh
cd ../../

cd ./reception/
govendor update +v
exec ./run.sh
cd ..
