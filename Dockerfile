FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache protobuf-dev git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

COPY internal/handler/schema.graphql /app/internal/handler/schema.graphql

RUN go mod tidy

RUN go install github.com/99designs/gqlgen@v0.17.73
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/order.proto

RUN $(go env GOPATH)/bin/gqlgen generate

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./main.go

FROM alpine:latest

RUN apk add --no-cache curl ca-certificates

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz && \
    mv migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

WORKDIR /root/

COPY --from=builder /order-service .
COPY --from=builder /app/migrations /migrations
COPY entrypoint.sh .

RUN chmod +x /root/entrypoint.sh

EXPOSE 8080 50051 8081

ENTRYPOINT ["/root/entrypoint.sh"]
