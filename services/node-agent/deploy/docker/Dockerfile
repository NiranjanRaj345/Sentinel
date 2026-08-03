FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o sentinel-agent ./cmd/sentinel-agent

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /app
COPY --from=builder /app/sentinel-agent .
COPY deploy/systemd/sentinel-agent.service /etc/systemd/system/
EXPOSE 8080
VOLUME ["/app/config.yaml", "/app/sentinel.db", "/app/logs"]
ENTRYPOINT ["/app/sentinel-agent"]
CMD ["run"]
