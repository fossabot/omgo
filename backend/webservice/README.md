Web Service
===

## Stages

### Register/Login (HTTPS)

acquire valid token

### Key Exchange (HTTPS)

a two step key exchange stage (Diffie-Hellman)

### Other API (HTTP)

invoke web service APIs via `GET` method

```json
{
"header": {
    "token": "base64 token",
    "param": "encrypted json string",
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
"signature": "sha1 checksum (key+value)"
}
```

**the calculation of signature**  

1. sort the keys of parameters in alphabet order and concat with value
2. replace special characters
3. calculate sha1 with shared secret

```
signature = sha1(replace(key1+value1+key2+value2+...keyn+valuen))
```

the signature will be verified after decrypt param and parse to json object
