# Dockerfile for custom github action
FROM golang:1.14 as build-img

COPY . $GOPATH/github.com/muncus/devto-publish-action
WORKDIR $GOPATH/github.com/muncus/devto-publish-action
RUN go get -d -v ./...
RUN go build -o /devto-sync ./cmd/devto-sync

FROM scratch
COPY --from=build-img /devto-sync /devto-sync
ENTRYPOINT ["sh", "-c", "/devto-sync", "$@"]
CMD ["--help"]