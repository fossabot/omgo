#!/bin/bash

cd ./ecdh
govendor update +v
cd ..

cd ./services
govendor update +v
cd ..

cd ./utils
govendor update +v
cd ..
