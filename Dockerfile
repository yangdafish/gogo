FROM golang:1.10

COPY . /code
WORKDIR /code

RUN go get -u github.com/gorilla/mux

RUN go build src/client.go
