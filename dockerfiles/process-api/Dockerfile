# syntax=docker/dockerfile:1

##
## BULD
##

FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache build-base git # ca-certificates tzdata

ARG SRC_DIR=./services/process-api
ARG PROGRAM_NAME=process-api

WORKDIR /src

ENV GOPATH=""
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY $SRC_DIR/go.mod $SRC_DIR/go.sum ./

COPY ./services/process-executor /deps/process-executor
COPY ./pkg /deps/lro
RUN go mod edit -replace github.com/lejenome/lro/services/process-executor=/deps/process-executor
RUN go mod edit -replace github.com/lejenome/lro/pkg=/deps/lro
RUN go mod edit -replace github.com/lejenome/lro/pkg=/deps/lro /deps/process-executor/go.mod

RUN go mod download && go mod verify
RUN go version

COPY $SRC_DIR .

RUN go mod edit -replace github.com/lejenome/lro/services/process-executor=/deps/process-executor
RUN go mod edit -replace github.com/lejenome/lro/pkg=/deps/lro

#RUN go vet -v
#RUN go test -v
RUN go build -v -trimpath -a -ldflags="-w -s" -o /app \
        ./cmd/$PROGRAM_NAME/main.go
# RUN ldd /app/bin/app

##
## Deploy
##

FROM scratch
#FROM gcr.io/distroless/static  # static go bin
#FROM gcr.io/distroless/base  # go bin with glibc
WORKDIR /

COPY --from=builder /app .

EXPOSE 8080
ENTRYPOINT ["/app"]
