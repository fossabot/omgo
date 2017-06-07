# Services

service management

## Usage

1. create a service root directory on ETCD

```
http://YOUR-ETCD-HOST/v2/keys/backends/
```

2. create sub directory for each kind of service

```
http://YOUR-ETCD-HOST/v2/keys/backends/agent
http://YOUR-ETCD-HOST/v2/keys/backends/game
http://YOUR-ETCD-HOST/v2/keys/backends/snowflake
```

3. add your service to these sub directories

```
curl -L -X PUT http://YOUR-ETCD-HOST/v2/keys/backends/agent/agent-001 -d value="127.0.0.1:27015"
```

Then you can init services

```go
services.Init("backends", ["127.0.0.1:2379"], ["snowflake", "game"])
```
