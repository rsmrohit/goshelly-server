# syntax=docker/dockerfile:1
FROM --platform=linux/amd64 ubuntu:latest  
WORKDIR /app
# COPY /scripts/certGen.sh ./
COPY *.sh ./
# COPY certGen.sh /
RUN chmod +x /certGen.sh && /certGen.sh



FROM --platform=linux/amd64 golang:1.16-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /goshelly-serv
EXPOSE 443
CMD [ "/goshelly-serv", "demo"] 