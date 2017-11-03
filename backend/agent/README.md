# agent

#### Packet Structure

```
|   2 bytes    | 4 bytes | varient  |
| payload size |   cmd   | protobuf |
```

```
2 bytes | package size
1 byte  | header size
1 byte  | protocol version
4 bytes | command
4 bytes | sequence
8 bytes | user serial number
```

```c
typedef struct PacketHeader {
    uint16_t pkgLen;  // 包大小(unsigned short)
    uint8_t headLen;  // 头部大小
    uint8_t version;  // 协议版本，当前为1
    int32_t cmd;  // 协议命令字
    int64_t uin;  // 用户账户uin（唯一的一个标识，跟登录相关）
    int32_t seq;  // 包的seq
} PacketHeader;
```

#### Pipeline  
 
PIPELINE #1 main.go  
PIPELINE #2 agent.go  
PIPELINE #3 buffer.go  

#### Key Exchange

