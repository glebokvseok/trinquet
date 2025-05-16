FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache curl git openssh protobuf gcc libc-dev

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN mkdir -p ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    chmod 644 ~/.ssh/known_hosts

COPY services/chat/go.mod services/chat/go.sum ./

RUN --mount=type=ssh go mod download

COPY services/chat/config ./config
COPY services/chat/internal ./internal

# Kafka Go client is based on the C library librdkafka
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -tags musl -o ./bin/main ./internal/cmd

CMD ["./bin/main"]
