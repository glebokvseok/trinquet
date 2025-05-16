FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache curl git openssh protobuf

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN mkdir -p ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    chmod 644 ~/.ssh/known_hosts

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY proto/ ./proto/

ENV GENERATED_PATH="/usr/local/go/src/generated"

RUN mkdir -p ${GENERATED_PATH} && \
    protoc \
        --go_out=${GENERATED_PATH} \
        --go-grpc_out=${GENERATED_PATH} \
        --proto_path=./proto \
         ./proto/**/*.proto

COPY services/player/go.mod services/player/go.sum ./

RUN --mount=type=ssh go mod download

COPY services/player/config ./config
COPY services/player/internal ./internal

RUN go build -o ./bin/main ./internal/cmd

CMD ["./bin/main"]
