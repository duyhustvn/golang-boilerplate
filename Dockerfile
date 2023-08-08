FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY cmd /src/cmd
COPY pkg /src/pkg
COPY internal /src/internal
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
RUN  go build -o main.o /src/cmd/service/main.go


FROM alpine:3.18.2
WORKDIR /src
RUN mkdir -p /src/logs
COPY --from=builder /src/main.o /src
CMD ["/src/main.o"]
