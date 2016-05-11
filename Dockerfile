FROM golang

ADD . /go/src/github.com/unitrans/unitrans
WORKDIR /go/src/github.com/unitrans/unitrans

ENV GO15VENDOREXPERIMENT 1
RUN go version && go env

RUN GO15VENDOREXPERIMENT=1 go install

CMD /go/bin/transpoint.io

EXPOSE 8088
