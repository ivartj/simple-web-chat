Simple Web Chat
===============

Simple WebSocket-based web chat application made in connection to a
job application.

To download and build requires [Go](https://golang.org/):

    $ go get github.com/ivartj/simple-web-server

To run:

    $ cd $GOPATH/src/github.com/ivartj/simple-web-chat
    $ simple-web-chat

The default assets directory is ./assets (relative to current working
directory). It can be adjusted with the --assets command-line option,
in which case the program does not need to be started in the source
directory.

To build docker image and run it as a container:

    $ docker build -t simple-web-chat .
    $ docker run -p 80:80 simple-web-chat

To stop docker container:

    $ docker stop $(docker ps -q -f ancestor=simple-web-chat)


