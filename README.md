vcenter-metrics
=================

[![Build Status](https://travis-ci.org/lenfree/vcenter-metrics.svg?branch=master)](https://travis-ci.org/lenfree/vcenter-metrics.svg?branch=master)

A bot that gather vcenter cluster and VMs summary and ships to Graphite and logging
backend such as Logstash. This is so we could measure, monitor, alert and build 
an event driven platform and decided to share with everyone

[Binary Releases](https://github.com/lenfree/vcenter-metrics/releases)

NOTE:

Since I moved this from a working repo which means I modified Dockerfile WORKDIR
to reflect new path and have yet to test build

## Usage

Copy .env.example to .env
$ cp .env.example .env

Uncompiled:
$ dotenv make run

Compiled:
$ dotenv ./vcenter-metrics

## Unit Test

### TODO

## Build

$ make release

## Contributing

1. Fork it ( https://github.com/lenfree/vcenter-metrics?fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
