FROM golang:1.23-alpine AS builder

WORKDIR /app
# ENV GOPROXY=https://goproxy.io,direct

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o build/baseTemp main.go


FROM debian:stable-slim
WORKDIR /app

VOLUME [ "/app/conf" ]
VOLUME [ "/app/logs" ]
EXPOSE 8080

COPY --from=builder /app/build/baseTemp .
COPY --from=builder /app/conf /app/conf


CMD ["/app/baseTemp","--conf=/app/conf","--port=8080"]