#!/bin/bash

cd ./utils/
rm -r vendor
govendor init
govendor add +external
cd ../

cd ./security/ecdh/
rm -r vendor
govendor init
govendor add +external
cd ../../

cd ./etcdclient/
rm -r vendor
govendor init
govendor add +external
cd ../

cd ./services/
rm -r vendor
govendor init
govendor add +external
cd ../

cd ./client/cli/
rm -r vendor
govendor init
govendor add +external
cd ../../
