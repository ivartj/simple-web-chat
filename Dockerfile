FROM golang:1.9

WORKDIR /go/src/github.com/ivartj/simple-web-chat
COPY ./*.go /go/src/github.com/ivartj/simple-web-chat/
COPY ./assets /go/src/github.com/ivartj/simple-web-chat/assets

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 80

CMD ["go-wrapper", "run" ]

