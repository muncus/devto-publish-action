# Dockerfile for custom github action
FROM golang:1.19

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o /devto-sync ./cmd/devto-sync

# ENTRYPOINT ["sh", "-c", "echo /devto-sync"]
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD []
