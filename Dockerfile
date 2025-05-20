FROM golang:1.23.3-bullseye AS builder

RUN apt-get update && \
    apt-get install -y git ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/${SERVICE}

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /bin/app /bin/app
COPY config.docker.json /app/config.json

ENV CONFIG_PATH=/app/config.json

ENTRYPOINT ["/bin/app"]
CMD ["--config", "/app/config.json"]
