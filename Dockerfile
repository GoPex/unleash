# From official golang alpine image
FROM gopex/ubuntu_golang:1.6rc1
MAINTAINER Albin Gilles "gilles.albin@gmail.com"
ENV REFRESHED_AT 2016-02-16

# Port exposed by this application
EXPOSE 3000

# Prepare directory holding our application
RUN mkdir -p $GOPATH/src/bitbucket.org/gopex/unleash

# Set the working directory
WORKDIR $GOPATH/src/bitbucket.org/gopex/unleash

# Copy our application into the container
COPY . .

# Get dependencies
RUN go get
