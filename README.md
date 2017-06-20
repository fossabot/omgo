# omgo
personal golang library

db -> reception -> auth -> agent -> game

### TODO

- [ ] game server
- [ ] service online/offline handling via ETCD
- [x] remove `omitempty` tag in .proto files
- [x] session management in agent
- [x] latch for mongodb in dbservice
- [x] add mechanism to agent for kick connections (gRPC ?)
- [x] rewrite cli client