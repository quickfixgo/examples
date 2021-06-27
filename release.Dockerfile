FROM golang:alpine
ADD config config
ADD qf /qf
ENTRYPOINT ["/qf"]