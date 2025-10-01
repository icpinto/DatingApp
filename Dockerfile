# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/dating-app ./cmd/app

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/dating-app /usr/local/bin/dating-app
COPY --from=builder /src/docs ./docs

RUN mkdir -p /app/uploads && chown -R app:app /app

USER app

EXPOSE 8080

ENTRYPOINT ["dating-app"]
