# syntax=docker/dockerfile:1
FROM alpine:latest 
# FROM darwinzeng/darwin-container:latest
WORKDIR /goshelly-server

COPY basic/* basic/
COPY bin/* bin/
COPY cmd/* cmd/
COPY scripts/* scripts/
COPY template/* template/
COPY *.mod /
COPY *.sum /
COPY *.go /


EXPOSE 443
RUN apk add --no-cache --upgrade bash
SHELL ["/bin/bash", "-c"] 
RUN apk add --update openssl && \ 
    rm -rf /var/cache/apk/*

#for linux image    
RUN chmod +x ./bin/app-amd64-linux
RUN chmod +x ./scripts/certGen.sh
RUN ls ./bin

# RUN ./bin/app-amd64-linux config
RUN chmod +x  ./scripts/goshelly-run-start.sh
CMD [ "./scripts/goshelly-run-start.sh"]


#for darwin image: BASE IMAGE DNE
# RUN chmod +x ./bin/app-amd64-darwin
# RUN [ "./bin/app-amd64-darwin", "config"]
# CMD [ "./bin/app-amd64-darwin", "demo"]
