# From official golang alpine image
FROM gopex/ubuntu_golang:1.6rc1
MAINTAINER Albin Gilles "albin.gilles@gmail.com"
ENV REFRESHED_AT 2016-02-04

# Port exposed by this application
EXPOSE 3000

# Prepare directory holding our application
RUN mkdir $GOPATH/src/github.com/GoPex/unleash

# Configure the GOPATH to our application
WORKDIR $GOPATH/src/github.com/GoPex/unleash
COPY . .
