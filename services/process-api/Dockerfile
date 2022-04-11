# syntax=docker/dockerfile:1

##
## BULD
##

FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache build-base git # ca-certificates tzdata

WORKDIR /src

ENV GOPATH=""
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go version

COPY . .
#RUN go vet -v
#RUN go test -v
RUN go build -v -trimpath -a -ldflags="-w -s" -o /app \
        ./cmd/process-api/main.go
# RUN ldd /app/bin/app

##
## Deploy
##

FROM scratch
#FROM gcr.io/distroless/static
WORKDIR /

COPY --from=builder /app .

EXPOSE 8080
ENTRYPOINT ["/app"]
