FROM golang:alpine

ADD config config

ADD bin/executor /executor

ENTRYPOINT ["/executor"]