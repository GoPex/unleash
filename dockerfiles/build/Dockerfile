FROM golang:1.6.0
MAINTAINER Albin Gilles <gilles.albin@gmail.com>
ENV REFRESHED_AT 2016-03-05

# Install upx-ucl utilities for compressing our go binary and git for go get
RUN apt-get update -yqq && apt-get install -yqq upx-ucl git

# Get the go goupx utility
RUN go get github.com/pwaller/goupx

# Copy the go package to build into the Docker Image
ENV PACKAGE github.com/GoPeX/unleash
ENV PACKAGE_DIR /go/src/$PACKAGE
RUN mkdir -p $PACKAGE_DIR
COPY . $PACKAGE_DIR/.

# Build the go package
RUN go get -t -d -v $PACKAGE/...
RUN CGO_ENABLED="${CGO_ENABLED:-0}" go build -a --installsuffix cgo --ldflags="${LDFLAGS:--s}" -o $GOPATH/bin/unleash $PACKAGE_DIR/main/server.go
RUN goupx $GOPATH/bin/unleash
