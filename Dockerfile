# Use golang alpine as base image
FROM golang:1.6.0-alpine
MAINTAINER Albin Gilles <gilles.albin@gmail.com>

# Add certificates
RUN apk add --update wget ca-certificates && \
  apk del wget ca-certificates && \
  rm /var/cache/apk/*

# Expose port to be used by unleash
EXPOSE 3000
