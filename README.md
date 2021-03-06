# authinator

This is a simple authentication server which stores users and issues tokens, it can either be used standalone or be embedded into an existing service.

# overview

* Uses [scrypt](www.tarsnap.com/scrypt.html) for password hashing
* Simple RESTish interface (see below)
* Issues [JWT](https://jwt.io/) tokens
* No Sessions or cookies

# REST API

TODO

## Register

```
curl -v -H "Content-Type: application/json; charset=UTF-8" -X POST \
  -d '{"login":"mememe","email":"me@example.com","password":"mepass"}' http://localhost:9090/users
```

## Login

```
curl -v --data "login=mememe&password=mepass" http://localhost:9090/auth/sign_in
```

## Update current user

```
curl -v -H "Content-Type: application/json; charset=UTF-8" \
  -H "Authorization: Bearer AS_ABOVE" \
  -X PUT -d '{"name":"Me Me"}' http://localhost:9090/users
```

## Get current User

```
curl -v -H "Content-Type: application/json; charset=UTF-8" \
  -H "Authorization: Bearer AS_ABOVE" \
  -X GET http://localhost:9090/users
```

# dependencies

* A data store, at the moment it supports RethinkDB with more to come.

# Features

* [x] Stores users
* [x] Authenticates users
* [x] Supports RethinkDB as a datastore
* [ ] Support for scopes and permission checks based on them
* [ ] Email activation of accounts
* [ ] Web interface

# references

* [JSON Web Token RFC 7519](https://tools.ietf.org/html/rfc7519)

# License

This project is released under BSD 3-clause license. Copyright 2016, [Mark Wolfe ](mailto:mark@wolfe.id.au).
