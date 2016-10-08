FROM golang:1.7.1
MAINTAINER lenfree.yeung@gmail.com

ENV GOPATH=/go

ENV LOGGER_HOST=logger
ENV LOGGER_PORT=4500
ENV METRIC_HOST=metrics
ENV METRIC_PORT=2003
ENV ENVIRONMENT=production

WORKDIR /go/src/github.com/lenfree/vcenter-metrics

RUN mkdir -p /go/src/github.com/lenfree/vcenter-metrics
COPY . /go/src/github.com/lenfree/vcenter-metrics/

CMD ["make start"]
