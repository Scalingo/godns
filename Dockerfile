FROM golang:1.20
LABEL maintainer="Étienne Michon <etienne@scalingo.com>"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/godns

EXPOSE 5321/udp

CMD $GOPATH/bin/godns
