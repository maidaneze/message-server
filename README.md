# Message-server

Backend server for messaging application written in golang that allows registration of users, 
sending and downloading text, voice and video messages to other users. 
The server uses basic user authentication via tokens. 
User information and messages are stored in an sqlite-3 database with the passwords encrypted.


## Contents

- [Running the server](#running-the-server)
- [Testing](#testing)
- [API Specification](#api-specification)
- [API Examples](#api-examples)
    - [Health Check](#health-check)
    - [Register User](#register-user)
    - [Login User](#login-user)
    - [Send Message](#send-message)
    - [Get Message](#get-message)

### Building the docker image

```
docker build --tag=message-server .
```

### Running the server

```
docker run -v $(pwd)/db:/server/db -p 8080:8080 message-server
```

### Testing

Requires go 1.11 and go modules.

```
go test ./...
```

### API Specification

The specification for the server is in the "swagger.json" file.

### API Examples

#### Health check

Check the health of the system.

Example:
```
curl -v GET localhost:8080/check
{"health":"ok"}
```

#### Register User

Registers a new user to the message service.
Requires an username that doesn't exist on the system and a non empty password.
On success it returns the userId of the new user.

Example:
```
curl -v POST -d '{"username":"testUser", "password":"testPassword"}' localhost:8080/users
{"id":1}
```

#### Login User

Login an user to the system. Requires the user's username and password.
On success it returns the userId of the user and an authorization token.

Example:
```
curl -v GET -d '{"username":"testUser", "password":"testPassword"}' localhost:8080/login
{"id":1,"token":"4974e5a5-c2ef-4562-be0d-54ab937b590b"}
```

#### Send Message

Sends a message from one user to another one. Requires the sender userId, the recipient userId and
a message content. Requires an access token for authentication.
It returns the messageId and a timestamp on success.

Example:

```
curl -v POST -d '{"sender":1, "recipient":2, "content":{"type":"text","text":"testMessage"}}' -H 'Authorization:Bearer 4974e5a5-c2ef-4562-be0d-54ab937b590b' localhost:8080/messages
{"id":1,"timestamp":"2019-04-12T12:52:45Z"}
```


Content types:
 - Text
 ``
 {"type":"text", "message":"someMessage"}
 ``
 - Image
 ``
 {"type":"image", "url":"www.example.com", "width":640, "height":480}
 ``
 - Video
 ``
 {"type":"video", "url":"www.example.com", "source":"youtube"}
 ``
 
Valid sources for video content are "youtube" and "vimeo".   
#### Get Message

Returns the messages for an user starting from a given messageId.
The number of messages can be limited by the "limit" parameter (default 100).
Requires an access token for authentication.

Example:
```
curl -v GET -H 'Authorization:Bearer e1289aba-acf2-4b24-a5a7-8559d9831c6b' "localhost:8080/messages?id=2&start=1&limit=100"
{"messages":[{"id":1,"timestamp":"2019-04-12T12:52:45.642356449Z","sender":1,"recipient":2,"content":{"type":"text","text":"testMessage"}}]}
```
