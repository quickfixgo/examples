FROM golang:alpine
ADD config config
ADD bin/qf /qf
ENTRYPOINT ["/qf"]