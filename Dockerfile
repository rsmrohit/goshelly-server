# syntax=docker/dockerfile:1
FROM golang:latest as layer1
WORKDIR /goshelly-server

COPY basic/* basic/
COPY bin/* bin/
COPY cmd/* cmd/
COPY scripts/* scripts/
COPY template/* template/
COPY *.mod .
COPY *.sum .
COPY *.go .
RUN GOOS=linux  go build  -o ./bin/app-amd64-linux .
RUN ls ./bin

FROM ubuntu:latest as base
# FROM darwinzeng/darwin-container:latest

COPY --from=layer1 /goshelly-server /goshelly-server
EXPOSE 443
WORKDIR /goshelly-server
RUN ls ./bin
# RUN apk add --no-cache --upgrade bash


azureuser@intern-sa:~/goshelly-server$ ls
Dockerfile  aws    certs        deployment.yaml  goshelly-server-api  scripts
LICENSE     basic  cmd          go.mod           logs                 template
README.md   bin    config.yaml  go.sum           main.go
azureuser@intern-sa:~/goshelly-server$ vim Dockerfile 



















# syntax=docker/dockerfile:1
FROM golang:latest as layer1
WORKDIR /goshelly-server

COPY basic/* basic/
COPY bin/* bin/
COPY cmd/* cmd/
COPY scripts/* scripts/
COPY template/* template/
COPY *.mod .
COPY *.sum .
COPY *.go .
RUN GOOS=linux  go build  -o ./bin/app-amd64-linux .
RUN ls ./bin

FROM ubuntu:latest as base
# FROM darwinzeng/darwin-container:latest

COPY --from=layer1 /goshelly-server /goshelly-server
EXPOSE 443
WORKDIR /goshelly-server
RUN ls ./bin
# RUN apk add --no-cache --upgrade bash
                                                              17,32         Top

