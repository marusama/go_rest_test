

Curl demo:
1. Register user:

curl -i -X POST -d '{"login":"admin","password":"123"}' -H "Content-Type:application/json" http://localhost:8080/register

You should get result:
HTTP/1.1 200 OK
...
{
    "status" : "OK"
}

2. Login using basic auth
"YWRtaW46MTIz" - base64 for "admin:123"

curl -i -X POST -H "Authorization: Basic YWRtaW46MTIz" http://localhost:8080/login

You will get an access token. The result should look like:
HTTP/1.1 200 OK
...
{
    "access_token" : "LcJWsc6Fu4bcWgd4ZkCUEQ==",
	"exp_time" : "2015-09-26T21:22:40.1764963+03:00"
}

3. Testing that we are authorized using this token:

curl -i -H "Authorization: Token LcJWsc6Fu4bcWgd4ZkCUEQ==" http://localhost:8080/auth_test

You should get result:
HTTP/1.1 200 OK
...
{
    "authed" : "admin"
}

If you are not authorized or token is expired:
HTTP/1.1 401 Unauthorized
...

4. Logging out:

curl -i -X POST -H "Authorization: Token LcJWsc6Fu4bcWgd4ZkCUEQ==" http://localhost:8080/logout

You should get result:
HTTP/1.1 200 OK
...
{
    "status" : "OK"
}
