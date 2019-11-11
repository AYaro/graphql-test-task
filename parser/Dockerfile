FROM golang:alpine

RUN apk update && apk add --no-cache git
ADD . /go/parser
WORKDIR /go/parser
RUN go get -d -v
RUN go build -o /go/bin/parser 

ENTRYPOINT [ "/go/bin/parser" ]

RUN chmod +x /go/bin/parser 