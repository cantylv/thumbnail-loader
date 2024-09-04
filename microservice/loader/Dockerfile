FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl gcc musl-dev

ARG MODULE_NAME=service
ENV CGO_ENABLED=1

WORKDIR /home/${MODULE_NAME}

COPY go.mod /home/${MODULE_NAME}/go.mod
COPY go.sum /home/${MODULE_NAME}/go.sum
COPY config  /home/${MODULE_NAME}/config
COPY microservice  /home/${MODULE_NAME}/microservice
COPY vendor  /home/${MODULE_NAME}/vendor
COPY services  /home/${MODULE_NAME}/services

WORKDIR /home/${MODULE_NAME}/microservice/loader
RUN go build -o main ./cmd/main/main.go

# Service
FROM alpine:latest as production
ARG MODULE_NAME=service

WORKDIR /home/${MODULE_NAME}

COPY --from=builder /home/${MODULE_NAME}/config/config.yaml config/config.yaml
COPY --from=builder /home/${MODULE_NAME}/services/sqlite/data/database.db services/sqlite/data/database.db
COPY --from=builder /home/${MODULE_NAME}/microservice/loader/main microservice/loader/cmd/main/main

WORKDIR /home/${MODULE_NAME}/microservice/loader/cmd/main
RUN chown root:root main

CMD ["./main"]