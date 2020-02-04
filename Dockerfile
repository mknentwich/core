FROM golang:alpine AS build

ENV GOOS=linux
# ENV GO111MODULE=on

WORKDIR /go/src/github.com/mknentwich/core
RUN apk add --no-cache \
    git \
    musl-dev \
    build-base
COPY . /go/src/github.com/mknentwich/core
RUN go get ./...
RUN go install -ldflags '-s -w' ./cmd/core

# ---

FROM alpine
COPY --from=build /go/bin/core /core
WORKDIR /data
VOLUME  /data
EXPOSE  9400
ENTRYPOINT [ "/core" ]
