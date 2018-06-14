# Build/Run the Client in Your Local Environment

This process assumes you already have Golang installed on your system

## Building The Client

Running the following command will produce a ```client``` executable in your local directory

    go build src/client.go

## Running The Client

You can use the executable by calling

    ./client

The client takes the following input flags

    --name (string)
    --port (int)
    --bootnodes (comma delimited string)

## Using The Client

Pass messages to connected nodes using the route ```/whisper```.
The request must include the following url parameters to be considered valid.

    name (string) (NOTE: the name of the node you want to send the message to)
    message (string) (NOTE: the message you want to send)


# Build/Run the Client Using a Dockerized Golang environment

Building the golang docker image

    docker build -t gogo .

Start the docker image locally

    docker run -itd --name gogo_con gogo

Then exec into the container

    docker exec -it gogo_con /bin/bash

and use the already built ```client``` executable, same as if you were running the client nativly


# Testing the Client

Go into the /src folder. Run the command

    go test
