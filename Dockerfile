# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.16-alpine


WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.txt ./
COPY *.go ./
COPY *.sh ./
COPY *.env ./
COPY instr.* ./
COPY /certs/* ./certs/
RUN go mod download

# RUN chmod +x certGen.sh  //initially for generating SSL certificates
# RUN ./certGen.sh         //

RUN go build -o /goshelly-serv

EXPOSE 443

CMD [ "/goshelly-serv", "-a"] 
# CMD ["/goshelly-serv", "-fe", "instr.txt"]