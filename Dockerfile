# Dockerfile for custom github action
FROM golang:1.14

COPY . $GOPATH/github.com/muncus/devto-publish-action
WORKDIR $GOPATH/github.com/muncus/devto-publish-action
RUN go get -d -v ./...
RUN go build -o /devto-sync ./cmd/devto-sync

ENTRYPOINT ["/devto-sync"]