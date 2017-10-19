Diffie-Hellman Key Exchange
===

Original Version : https://github.com/wsddn/go-ecdh

Reference:
1. https://cr.yp.to/ecdh.html
2. https://en.wikipedia.org/wiki/Diffie%E2%80%93Hellman_key_exchange
3. https://code.google.com/archive/p/curve25519-donna/


```
Alice generates private key xA and public key eA
Bob generates private key xB and public key eB  
Alice send her public key to Bob
Bob send his public key to Alice
Alice create common secret = f(xA, eB)
Bob create common secret = f(xB, eA)
Alice now talks with Bob in channel encrypted by common secret
```