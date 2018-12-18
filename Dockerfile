FROM golang:latest

MAINTAINER Kohei Hosokawa

RUN apt-get update && apt-get install -y git

RUN go get -u golang.org/x/vgo

RUN mkdir -p /pac

COPY ./Makefile /pac/

WORKDIR /pac
