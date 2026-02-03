FROM golang:1.25.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o refrain ./cmd/refrain/main.go

FROM alpine:latest
WORKDIR /app
VOLUME [ "/data" ]
COPY --from=builder /app/refrain /app/refrain
RUN chmod +x /app/refrain
ENTRYPOINT ["/app/refrain"]
