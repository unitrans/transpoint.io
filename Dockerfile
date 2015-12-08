FROM golang

ADD . /go/src/github.com/urakozz/transpoint.io
WORKDIR /go/src/github.com/urakozz/transpoint.io

ENV HOME /root
ENV GOPATH /go
ENV GO15VENDOREXPERIMENT 1

RUN go install

CMD /go/bin/transpoint.io

EXPOSE 8088
