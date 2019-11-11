FROM golang:alpine

RUN apk update && apk add --no-cache git
RUN apk add --update go git build-base
ADD . /go/
WORKDIR /go/server
RUN go get -d -v
RUN go build -o /go/bin/graphql-test-task

ENTRYPOINT [ "/go/bin/graphql-test-task" ]

