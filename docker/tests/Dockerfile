FROM golang:latest
ADD . /turl
WORKDIR /turl
CMD CGO_ENABLED=0 go test ./...