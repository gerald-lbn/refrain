FROM golang:1.25.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o refrain ./cmd/refrain/main.go

FROM alpine:latest
WORKDIR /app
VOLUME [ "/data", "/config" ]
RUN apk add --no-cache su-exec shadow
RUN groupadd -g 1000 refrain && \
    useradd -u 1000 -g 1000 -d /app -s /bin/sh refrain
COPY --from=builder /app/refrain /app/refrain
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/refrain /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/refrain"]
