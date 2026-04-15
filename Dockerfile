FROM golang:1.23-alpine AS builder

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

WORKDIR /root/

COPY --from=builder /order-service .
COPY --from=builder /app/migrations /migrations

EXPOSE 8080 50051 8081

CMD ["./order-service"]
