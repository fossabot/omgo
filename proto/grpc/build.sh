#!/bin/bash

for f in $(pwd)/*;
do
  [ -d ${f} ] && cd "${f}" && protoc --go_out=plugins=grpc:. *.proto
done;

