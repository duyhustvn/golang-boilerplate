FROM golang:1.21-alpine AS builder

#ARG ARG_HTTP_PROXY 127.0.0.1:1234
#ARG ARG_HTTPS_PROXY 127.0.0.1:1234
#

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/service/main.go

RUN  go build -o main.o /src/cmd/service/main.go

FROM alpine:3.18.2 AS production-stage
WORKDIR /src
RUN mkdir -p /src/logs
COPY --from=builder /src/main.o ./golang-boilerplate-be
COPY --from=builder /src/docs/ ./docs/
CMD ["./golang-boilerplate-be"]
