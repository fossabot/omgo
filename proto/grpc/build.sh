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
  # remove unnecessary `omitempty` for bool key-pair
  outfile=${output}/${name}
  sed -e '/varint/ s/,omitempty//' ${outfile}.pb.go > ${outfile}.tmp
  mv ${outfile}.tmp ${outfile}.pb.go
done;

