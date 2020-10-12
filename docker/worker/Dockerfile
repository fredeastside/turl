# modules
FROM golang:latest as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download

# binary
FROM golang:latest as builder
COPY --from=modules /go/pkg /go/pkg
RUN useradd -u 10001 turl
RUN mkdir -p /turl
ADD . /turl
WORKDIR /turl
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/turl ./cmd/worker

# run
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER turl
COPY --from=builder /turl/bin/turl /turl
CMD ["/turl"]