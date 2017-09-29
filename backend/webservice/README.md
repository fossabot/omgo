Web Service
===

## Stages

### Register/Login (HTTPS)

acquire valid token

### Key Exchange (HTTPS)

key exchange stage (Diffie-Hellman)

### Other API (HTTP)

invoke web service APIs via `GET` method

```json
{
"header": {
    "token": "base64 token",
    "param": "encrypted json string",
    "signature": "sha1 checksum (key+value)",
    "nonce": "timestamp + request sequence"
}
}
```

the param is encrypted from parameter json object


```json
{
"key1": "value1",
"key2": "value2",
"key3": "value3",
"...": "value...",
"timestamp": "value"
}
```

the param is encrypt via shared secret

```
header.param = encrypt(paramJson.encode)
```

XXTEA will be good enough

**the calculation of signature**  

1. sort the keys of parameters in alphabet order and concat with value
2. replace special characters
3. add salt
4. calculate sha1 checksum

```
signature = sha1(salt(replace(key1+value1+key2+value2+...keyn+valuen)))
```

the signature will be verified after decrypt param and parse to json object
