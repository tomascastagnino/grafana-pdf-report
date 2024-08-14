FROM golang:1.20-alpine AS builder

WORKDIR /app/src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/server ./cmd/server
COPY internal ./internal
RUN go build -o /go/bin/grafana-pdf-reporter ./cmd/server

FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install

COPY static ./static

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app/src/cmd/server

COPY --from=builder /go/bin/grafana-pdf-reporter .
COPY --from=frontend-builder /app/static ../../static
COPY --from=frontend-builder /app/node_modules ../../node_modules
COPY config.ini ../../config.ini

EXPOSE 9090

CMD ["./grafana-pdf-reporter"]
