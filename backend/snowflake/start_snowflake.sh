#!/usr/bin/env bash
docker rm -f snowflake1
docker build --no-cache --rm=true -t snowflake .
docker run --rm=true --name snowflake1 -e SERVICE_ID=snowflake1 -e MACHINE_ID=1 --net=host -p 40001:40001 -d -P snowflake \
    -p 40001 \
    -e http://127.0.0.1:2379

curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/test_key -d value="0"
curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/userid -d value="0"
curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/snowflake-uuid -d value="0"
