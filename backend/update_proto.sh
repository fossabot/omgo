#!/usr/bin/env bash

# dbservice
cp ../proto/grpc/db.proto db/dbservice/src/main/proto
cp ../proto/grpc/snowflake.proto db/dbservice/src/main/proto

# webservice
cp ../proto/grpc/db.proto webservice/src/main/proto
