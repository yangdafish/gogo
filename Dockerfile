FROM golang:1.10

VOLUME /code

COPY . /code
WORKDIR /code

