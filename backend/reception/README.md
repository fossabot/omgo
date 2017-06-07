Auth
===

## Authentication module

port 40000 is for internal service  
port 8080 is for http login, and it should be changed to 80 or 443 for better accessibility

### Login with Email and password

email and password digest send via http header  

`password digest = hash(password + salt)`
