FROM golang:1.25.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o refrain ./cmd/refrain/main.go

FROM alpine:latest
WORKDIR /app
VOLUME [ "/music", "/config" ]
RUN apk add --no-cache su-exec && \
    addgroup -S refrain -g 1000 && \
    adduser -S refrain -G refrain -u 1000
COPY --chown=refrain:refrain --from=builder /app/refrain /app/refrain
COPY --chown=refrain:refrain entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/refrain /app/entrypoint.sh && \
    mkdir /music /config && \
    chown refrain:refrain /music /config
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/refrain"]
