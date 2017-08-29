#!/bin/bash

cd ./utils/
govendor update +v
cd ../

cd ./security/ecdh/
govendor update +v
cd ../../

cd ./services/
govendor update +v
cd ../

cd ./client/cli/
govendor update +v
cd ../../
