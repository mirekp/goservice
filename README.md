# goservice

A simple restful messaging service demo in Go. Features mux routing and JSON API interface.


Create user account:
```
curl -H "Content-Type: application/json" -X POST -d '{ "data" : {"id":"account1","type":"account","attributes":{"password":"xyza"}}}' http://localhost:9999/accounts
```
Post message:
```
curl -H "Content-Type: application/json" -X POST -d 'data : {"type": "message", "attributes": { "text", "this is a message text”} }' http://localhost:9999/accounts
```
Delete message:
```
curl -X DELETE http://localhost:9999/messages/1
```
Viewing messages:
```
curl http://localhost:9999/messages
```
Editing message:
```
curl -H "Content-Type: application/json" -X PATCH -d '{"data" : {"attributes": { "text": "this is a new message text"}}}' http://localhost:9999/messages/1
```
