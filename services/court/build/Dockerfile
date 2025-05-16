FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache curl git openssh

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN mkdir -p ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    chmod 644 ~/.ssh/known_hosts

COPY services/court/go.mod services/court/go.sum ./

RUN --mount=type=ssh go mod download

COPY services/court/config ./config
COPY services/court/internal ./internal

RUN go build -o ./bin/main ./internal/cmd

CMD ["./bin/main"]
