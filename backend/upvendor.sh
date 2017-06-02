#!/bin/bash

cd ./snowflake/
govendor update +v
cd ..

cd ./agent/
govendor update +v
cd ..

cd ./db/dbservice/
govendor update +v
cd ../../

cd ./reception/
govendor update +v
cd ..
