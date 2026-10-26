FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gotraker ./cmd/api/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/gotraker .
EXPOSE 8080
CMD ["./gotraker"]