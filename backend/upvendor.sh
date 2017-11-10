#!/bin/bash

cd ./snowflake/
dep ensure -update
dep prune
cd ..

cd ./agent/
dep ensure -update
dep prune
cd ..

cd ./game/
dep ensure -update
dep prune
cd ..
