# agent

#### Packet Structure

client -> server

```
2 bytes | header size (these 2 bytes excluded)
n bytes | header protobuf
n bytes | payload protobuf
```

server -> client

```
2 bytes | header size (these 2 bytes excluded)
4 bytes | cmd
n bytes | payload
```

#### Pipeline  
 
PIPELINE #1 main.go  
PIPELINE #2 agent.go  
PIPELINE #3 buffer.go  

#### Key Exchange

