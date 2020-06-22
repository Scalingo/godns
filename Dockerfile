FROM golang:1.14
MAINTAINER Étienne Michon "etienne@scalingo.com"

RUN go get github.com/cespare/reflex

WORKDIR $GOPATH/src/github.com/Scalingo/godns
COPY . ./

RUN go install

EXPOSE 5353

CMD $GOPATH/bin/godns
