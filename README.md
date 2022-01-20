## Test Setup
```
make run
```
will setup a test environment to test the service

## Endpoint Usage
### Headers
Set authentication request headers:
```
X-KBU-Auth: abcdefghijklmnopqrstuvwxyz
```

### Requests
#### GET
```
/ping | responds 200 OK
```

#### POST
Registration:
```
/register | responds 201 CREATED if successfull
Payload:
{
    "user": "USERNAME",
    "password": "PASSWORD",
    "role": "ROLE",
    "email": "EMAIL"
}
```
Login:
```
/login | responds 200 OK if logged in
Payload:
{
    "password": "PASSWORD",
    "email": "EMAIL"
}

Response:
{
    "token": "PASETO TOKEN"
}
```
Token validation:
```
/validate | responds 200 OK if paseto token valid, 401 if not valid
Payload:
{
    "token": "PASETO TOKEN"
}
Response:
{
    "role": "USER_ROLE"
}
```
Password Reset Request:
```
/resetpassword | responds 200 OK if successfull
Payload:
{
    "email": "EMAIL"
}

Response:
{
    "resetId": "PASSWORD_RESET_ID"
}
```
Password Reset:
```
/reset?resetId="PASSWORD_RESET_ID" | Responds 200 OK if successfull
Payload:
{
    "password": "PASSWORD"
}
```
