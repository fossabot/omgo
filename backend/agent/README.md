# agent

#### Packet Structure

client -> server

```
2 bytes | header size (these 2 bytes excluded)
n bytes | header protobuf (unencrypted)
n bytes | payload protobuf (might be encrypted)
```

server -> client

```
2 bytes | header size (these 2 bytes excluded)
n bytes | response header protobuf
n bytes | payload
```

#### Pipeline  
 
PIPELINE #1 main.go  
PIPELINE #2 agent.go  
PIPELINE #3 buffer.go  

#### Key Exchange

