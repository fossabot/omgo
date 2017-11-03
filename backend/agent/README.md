# agent

#### Packet Structure

```
|   2 bytes    | 4 bytes | varient  |
| payload size |   cmd   | protobuf |
```

```
2 bytes | header size (these 2 bytes excluded)
n byte  | header protobuf
```

#### Pipeline  
 
PIPELINE #1 main.go  
PIPELINE #2 agent.go  
PIPELINE #3 buffer.go  

#### Key Exchange

