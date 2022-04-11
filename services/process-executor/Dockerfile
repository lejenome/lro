# syntax=docker/dockerfile:1

##
## BULD
##

FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache build-base git # ca-certificates tzdata

WORKDIR /app/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -a -v -o /app/bin/app ./cmd/process-executor/main.go
RUN go build -v -o /app/bin/app ./cmd/process-executor/main.go
RUN ldd /app/bin/app

##
## Deploy
##

FROM scratch
WORKDIR /

COPY --from=builder /app/bin/app .

EXPOSE 8080
ENTRYPOINT ["/app"]
