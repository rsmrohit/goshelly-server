# syntax=docker/dockerfile:1
# FROM --platform=linux/amd64 alpine:latest
# WORKDIR /app
# COPY . .
# EXPOSE 443
# EXPOSE 8080
# RUN chmod +x ./scripts/certGen.sh
# RUN ./scripts/certGen.sh
# # RUN chmod +x ./scripts/goshelly-run-start.sh



FROM golang:1.16-alpine

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /goshelly-serv

EXPOSE 443
RUN /goshelly-serv config
CMD [ "/goshelly-serv", "demo" ]
