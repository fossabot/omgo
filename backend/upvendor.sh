#!/bin/bash

cd ./snowflake/
govendor update +v
cd ..

cd ./agent/
govendor update +v
cd ..

cd ./game/
govendor update +v
cd ..
