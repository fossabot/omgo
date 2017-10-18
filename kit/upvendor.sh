#!/bin/bash

cd ./ecdh
govendor update +v
cd ..

cd ./services
govendor update +v
cd ..

cd ./util
govendor update +v
cd ..
