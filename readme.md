# go-rest-test

Simple REST API demo server written on go (golang) uses Cassandra as storage.
Demonstrates base authorization and authentication with token.

A client can register new user, log in, call authorization required test function, and log out.

## Used libraries:
- github.com/ant0ine/go-json-rest - A quick and easy way to setup a RESTful JSON API.
- github.com/grayj/go-json-rest-middleware-tokenauth - Token authentication middleware for go-json-rest.
- https://github.com/gocql/gocql - Package gocql implements a fast and robust Cassandra client for the Go programming language. 

## API methods:
#### 1. http://localhost:8080/register
*Request:*

HTTP Method: POST. No authorization required. 

Request payload: JSON
```JSON
{
    "login" : "new_login_here",
    "password" : "new_password_here"
}
```

*Response:*
- HTTP 200: OK. User registered.

Body: ```{ "status" : "OK" } ```
- HTTP 400: Bad request. Invalid input JSON; Login is empty; Password is empty; User is already registered.

Body: ``` { "Error" : "Error message here" } ```

#### 2. http://localhost:8080/login
*Request:*

HTTP Method: POST. Basic access authentication required (https://en.wikipedia.org/wiki/Basic_access_authentication). Use early registered user login/password.

Request payload: none.

*Response:*
- HTTP 200: OK. User logged in.

Body (sample):
```
{
    "access_token" : "LcJWsc6Fu4bcWgd4ZkCUEQ==",
	"exp_time" : "2015-09-26T21:22:40.1764963+03:00"
}
```
You will get a token and datetime when the token will expire.
- HTTP 400: Bad request. Invalid authentication.

Body: ``` { "Error" : "Invalid authentication" } ```
#### 3. http://localhost:8080/auth_test
*Request:*

HTTP Method: GET. Token access authentication required (http://www.w3.org/2001/sw/Europe/events/foaf-galway/papers/fp/token_based_authentication/). Use early obtained token.

Request payload: none.

*Response:*
- HTTP 200: OK. User is authenticated.

Body: ``` {"authed" : "some username here" } ```
You will get a login for provided token.
- HTTP 401: Unauthorized. Invalid authentication token or token expired.

Body: ``` { "Error" : "Not authorized" } ```
#### 4. http://localhost:8080/logout
*Request:*

HTTP Method: POST. Token access authentication required. Use early obtained token.

Request payload: none.

*Response:*
- HTTP 200: OK. User is logged out.

Body: ``` {"status" : "OK"} ```
- HTTP 401: Unauthorized. Invalid authentication token or token expired.

Body: ``` { "Error" : "Not authorized" } ```

## Curl demo:
#### 1. Register user:

```sh
curl -i -X POST -d '{"login":"admin","password":"123"}' -H "Content-Type:application/json" http://localhost:8080/register
```

You should get result:

```json
HTTP/1.1 200 OK
...
{
    "status" : "OK"
}
```

#### 2. Login using basic auth

"YWRtaW46MTIz" - base64 for "admin:123"

```sh
curl -i -X POST -H "Authorization: Basic YWRtaW46MTIz" http://localhost:8080/login
```

You will get an access token. The result should look like:

```json
HTTP/1.1 200 OK
...
{
    "access_token" : "LcJWsc6Fu4bcWgd4ZkCUEQ==",
	"exp_time" : "2015-09-26T21:22:40.1764963+03:00"
}
```

#### 3. Testing that we are authorized using this token:

```sh
curl -i -H "Authorization: Token LcJWsc6Fu4bcWgd4ZkCUEQ==" http://localhost:8080/auth_test
```

You should get result:

```json
HTTP/1.1 200 OK
...
{
    "authed" : "admin"
}
```

If you are not authorized or token is expired the response will be:
```json
HTTP/1.1 401 Unauthorized
...
```

#### 4. Logging out:

```sh
curl -i -X POST -H "Authorization: Token LcJWsc6Fu4bcWgd4ZkCUEQ==" http://localhost:8080/logout
```

You should get result:
```json
HTTP/1.1 200 OK
...
{
    "status" : "OK"
}
```
