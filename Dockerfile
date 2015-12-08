FROM golang

ADD . /go/src/github.com/urakozz/transpoint.io
WORKDIR /go/src/github.com/urakozz/transpoint.io

ENV GO15VENDOREXPERIMENT 1

RUN go install

CMD /go/bin/transpoint.io

EXPOSE 8088