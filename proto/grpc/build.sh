#!/bin/bash

# iterate .proto files
for f in $(pwd)/*.proto;
do
  # get proto file name as output directory name
  name=$(basename "$f" ".proto")
  # make directory
  output=$(pwd)/${name}
  mkdir -p ${output}
  # compile
  parent=$(dirname "$(pwd)")
  protoc -I . -I ../ --go_out=plugins=grpc:${name} ${name}.proto 

  #[ -d ${f} ] && cd "${f}" && protoc -I=../ --go_out=plugins=grpc:. *.proto
done;

