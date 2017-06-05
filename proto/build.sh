#!/bin/bash

# protoc -I=. --go_out=./pb *.proto

# make output root directory
mkdir -p pb

# iterate .proto files
for f in $(pwd)/*.proto;
do
  # get proto file name as output directory name
  name=$(basename "$f" ".proto")
  # make directory
  output=$(pwd)/pb/${name}
  mkdir -p ${output}
  # compile
  protoc -I=. --go_out=${GOPATH}/src ${name}.proto
  # remove unnecessary `omitempty` for bool key-pair
  outfile=${output}/${name}
  sed -e '/varint/ s/,omitempty//' ${outfile}.pb.go > ${outfile}.tmp
  mv ${outfile}.tmp ${outfile}.pb.go
done;
