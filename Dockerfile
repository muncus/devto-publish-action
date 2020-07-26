# Dockerfile for custom github action
FROM golang:1.14

COPY . $GOPATH/github.com/muncus/devto-publish-action
RUN go get -d -v github.com/muncus/devto-publish-action/...
RUN go build -o /devto-sync github.com/muncus/devto-publish-action/cmd/devto-sync

# ENTRYPOINT ["sh", "-c", "echo /devto-sync"]
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD []