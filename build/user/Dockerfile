FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/user

FROM alpine:latest
WORKDIR /app
COPY --from=builder /usr/local/bin/app .
CMD ["./app"]