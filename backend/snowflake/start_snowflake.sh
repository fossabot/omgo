#!/usr/bin/env bash
SID=snowflake1
HOST=127.0.0.1

docker rm -f ${SID}
docker build --no-cache --rm=true -t snowflake .
docker run --rm=true --name ${SID} -e SERVICE_ID=${SID} -e MACHINE_ID=1 --net=host -p 40001:40001 -d -P snowflake \
    -p 40001 \
    -e http://${HOST}:2379

# register service
curl -L -X PUT http://${HOST}:2379/v2/keys/backends/snowflake/${SID} -d value=${HOST}:40001

# init etcd variables
curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/test_key -d value="0"
curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/userid -d value="0"
curl -L -X PUT http://127.0.0.1:2379/v2/keys/seqs/snowflake-uuid -d value="0"
