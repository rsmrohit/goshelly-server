# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /goshelly-serv
EXPOSE 443
SHELL ["/bin/bash", "-c"]
RUN ["/goshelly-serv" ,"config"]
CMD [ "/goshelly-serv", "demo" ]
