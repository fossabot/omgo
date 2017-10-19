# Services

service management

## Usage

1. setup ETCD v3 environment

2. add services to ETCD like

```plain
Key: {service-root}/{service-type}/{service-id}
Value: 127.0.0.1:40001
```

Then you can init services

```go
services.Init("backends", ["127.0.0.1:2379"], ["snowflake", "game"])
```

## TODO

add circuit breaker
add unregister service
