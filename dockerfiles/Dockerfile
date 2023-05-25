FROM golang:1.19-alpine3.17

COPY . /go/src/qubic-prime-bind-server

WORKDIR /go/src/qubic-prime-bind-server

RUN  go build -o build/local_web_server local_web_server/main.go