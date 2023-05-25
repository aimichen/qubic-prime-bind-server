# Qubic Third-Party Prime Bind Backend Example
This is sample code of 3rd-party prime bind backend, and [docker image](https://hub.docker.com/r/aimi/qubic-prime-bind-server) is ready for local testing.

## Run local server
```
$ docker pull aimi/qubic-prime-bind-server:latest
$ docker run -it --rm -p8080:80 aimi/qubic-prime-bind-server:latest build/local_web_server \
	--adminAPIKey={YOUR_API_KEY} \
	--adminAPISecret={YOUR_API_SECRET}
```

## Prime bind api
```
$ curl -X POST --data "bindTicket=mock-bind-ticket&memberId=999" localhost:8080/primeBind
{"success": true}

$ curl -X POST --data "memberId=999" localhost:8080/credentialIssue
{"identityTicket": "mock identity ticket", "expiredAt": "2023-05-25T15:22:58+08:00", "address": "0x1234"}%
```