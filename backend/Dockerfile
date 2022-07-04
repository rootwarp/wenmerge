# Build stage.
FROM golang:1.18.3-stretch AS builder

WORKDIR /app
ADD . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLE=0 go build -ldflags '-extldflags "-static"'

# Exec. stage.
FROM alpine:3.16.0

WORKDIR /app
COPY --from=builder /app/wenmerge  ./

EXPOSE 9090
CMD ["/app/wenmerge"]
