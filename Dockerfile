FROM golang:1.23.2 AS builder
WORKDIR /app
RUN go mod init utopia-client
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/utopia-client

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/utopia-client .
CMD ["./utopia-client"]
