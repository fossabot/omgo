#!/usr/bin/env bash

# dbservice
cp ../proto/grpc/db.proto dbservice/src/main/proto
cp ../proto/grpc/snowflake.proto dbservice/src/main/proto

# dataservice
cp ../proto/grpc/db.proto dataservice/src/main/proto
cp ../proto/grpc/snowflake.proto dataservice/src/main/proto

# webservice
cp ../proto/grpc/db.proto webservice/src/main/proto
