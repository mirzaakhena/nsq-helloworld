#!/bin/bash

rm -rf dist/mirserver-api
rm -rf dist/mirserver-con

GOOS=linux GOARCH=amd64 go build -o dist/mirserver-api app/api/main.go
GOOS=linux GOARCH=amd64 go build -o dist/mirserver-con app/consumer/main.go

docker build -f DockerfileApi -t mirserverapi .
docker build -f DockerfileCon -t mirservercon .
