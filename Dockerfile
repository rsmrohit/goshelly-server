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
COPY *.yaml .
COPY html/* html/
COPY Dockerfile .

COPY goshelly-server-api/* goshelly-server-api/
RUN GOOS=linux  go build  -o ./bin/app-amd64-linux .
RUN ls ./bin

FROM ubuntu:latest as base


COPY --from=layer1 /goshelly-server /goshelly-server
EXPOSE 443
EXPOSE 8000
EXPOSE 9000

WORKDIR /goshelly-server
RUN ls ./bin
# RUN apk add --no-cache --upgrade bash
RUN apt-get install --only-upgrade bash

SHELL ["/bin/bash", "-c"]
RUN apt-get update -y
RUN apt-get install openssl -y

#for linux image    
RUN chmod +x ./bin/app-amd64-linux
RUN chmod +x ./scripts/certGen.sh
RUN ls -altr ./bin

# RUN ./bin/app-amd64-linux config
RUN chmod +x  ./scripts/goshelly-run-start.sh
CMD [ "./scripts/goshelly-run-start.sh"]