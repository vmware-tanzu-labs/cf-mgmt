FROM golang:latest

MAINTAINER Caleb Washburn "cwashburn@pivotal.io"

ADD cf-mgmt /usr/bin/cf-mgmt
RUN /usr/bin/cf-mgmt
