# goservice

A simple restful messaging service demo in Go.

Features JSON API interface.

Prerequisites:

- go compiler

Dependencies

   - gorilla/mux
   - crypto/bcrypt


To create user account:
```
curl -H "Content-Type: application/json" -X POST -d '{ "data" : { "type": "account", "attributes": { "username" : "bob", "password": "password123"}}}' http://localhost:9999/accounts
```
To post a message:
```
curl -H "Content-Type: application/json" -X POST -d '{ "data": { "type": "message", "attributes": { "text": "this is a message text"}}}' http://localhost:9999/messages
```
Delete the message:
```
curl -X DELETE http://localhost:9999/messages/0
```
Viewing all messages:
```
curl http://localhost:9999/messages
```
Editing a message:
```
curl -H "Content-Type: application/json" -X PATCH -d '{ "data": {"attributes": { "text": "this is a new message text"}}}' http://localhost:9999/messages/0
```
