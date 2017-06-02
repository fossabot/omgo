#!/bin/bash

cd ./client/cli/
govendor update +v
cd ../../

cd ./etcdclient/
govendor update +v
cd ../

cd ./security/ecdh/
govendor update +v
cd ../../

cd ./services/
govendor update +v
cd ../

cd ./utils/
govendor update +v
cd ../
