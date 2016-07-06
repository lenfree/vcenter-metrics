vcenter-metrics
=================

A bot that gather vcenter cluster and VMs summary and ships to Graphite and logging
backend such as Logstash. This is so we could measure, monitor, alert and build 
an event driven platform and decided to share with everyone

NOTE:

Since I moved this from a working repo which means I modified Dockerfile WORKDIR
to reflect new path and have yet to test build

## Usage

Copy .env.example to .env
$ cp .env.example .env

Uncompiled:
$ dotenv make

Compiled:
$ dotenv ./vcenter-metrics

## Unit Test

### TODO

## Build

$ go build ./...

## Contributing

1. Fork it ( https://git02.ae.sda.corp.telstra.com/projects/TOP/repos/vcenter-metrics?fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
