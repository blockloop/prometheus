# Prometheus 

This is a fork of Prometheus that enables remote_write. It is used by @digitalocean/insights.

[![Build Status](https://travis-ci.org/prometheus/prometheus.svg)][travis]
[![CircleCI](https://circleci.com/gh/prometheus/prometheus/tree/master.svg?style=shield)][circleci]
[![Docker Repository on Quay](https://quay.io/repository/prometheus/prometheus/status)][quay]
[![Docker Pulls](https://img.shields.io/docker/pulls/prom/prometheus.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/prometheus/prometheus)](https://goreportcard.com/report/github.com/prometheus/prometheus)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/486/badge)](https://bestpractices.coreinfrastructure.org/projects/486)
[![fuzzit](https://app.fuzzit.dev/badge?org_id=prometheus&branch=master)](https://fuzzit.dev)
>>>>>>> upstream/master

The master branch is kept up to date with the latest stable tag of upstream prometheus. All fork changes are written to
the `remote_write` branch. 

To make changes or fixes to remote_write capability, push a PR to the remote_write branch.

To update the master branch from upstream make sure you have the upstream remote added locally

```
git remote add upstream git@github.com:prometheus/prometheus.git
```

Once you have the remote you can pull from upstream, checkout the latest stable tag, and then rebase it into
the remote_write branch

```
git checkout master
git pull upstream master
git reset --hard <latest stable tag, i.e. v2.8.1>
git checkout remote_write
git rebase master
```

Once you have rebased master, build it with docker


```
CGO_ENABLED=0 go build ./cmd/prometheus
docker build . -t digitalocean/prometheus:<latest stable tag, i.e. v2.8.1>
docker push digitalocean/prometheus:<tag>
```
There are various ways of installing Prometheus.

### Precompiled binaries

Precompiled binaries for released versions are available in the
[*download* section](https://prometheus.io/download/)
on [prometheus.io](https://prometheus.io). Using the latest production release binary
is the recommended way of installing Prometheus.
See the [Installing](https://prometheus.io/docs/introduction/install/)
chapter in the documentation for all the details.

Debian packages [are available](https://packages.debian.org/sid/net/prometheus).

### Docker images

Docker images are available on [Quay.io](https://quay.io/repository/prometheus/prometheus) or [Docker Hub](https://hub.docker.com/r/prom/prometheus/).

You can launch a Prometheus container for trying it out with

    $ docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus

Prometheus will now be reachable at http://localhost:9090/.

### Building from source

To build Prometheus from the source code yourself you need to have a working
Go environment with [version 1.12 or greater installed](https://golang.org/doc/install).

You can directly use the `go` tool to download and install the `prometheus`
and `promtool` binaries into your `GOPATH`:

    $ go get github.com/prometheus/prometheus/cmd/...
    $ prometheus --config.file=your_config.yml

You can also clone the repository yourself and build using `make`:

    $ mkdir -p $GOPATH/src/github.com/prometheus
    $ cd $GOPATH/src/github.com/prometheus
    $ git clone https://github.com/prometheus/prometheus.git
    $ cd prometheus
    $ make build
    $ ./prometheus --config.file=your_config.yml

The Makefile provides several targets:

  * *build*: build the `prometheus` and `promtool` binaries
  * *test*: run the tests
  * *test-short*: run the short tests
  * *format*: format the source code
  * *vet*: check the source code for common errors
  * *assets*: rebuild the static assets
  * *docker*: build a docker container for the current `HEAD`

## More information

  * The source code is periodically indexed: [Prometheus Core](https://godoc.org/github.com/prometheus/prometheus).
  * You will find a Travis CI configuration in `.travis.yml`.
  * See the [Community page](https://prometheus.io/community) for how to reach the Prometheus developers and users on various communication channels.

## Contributing

Refer to [CONTRIBUTING.md](https://github.com/prometheus/prometheus/blob/master/CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](https://github.com/prometheus/prometheus/blob/master/LICENSE).


[travis]: https://travis-ci.org/prometheus/prometheus
[hub]: https://hub.docker.com/r/prom/prometheus/
[circleci]: https://circleci.com/gh/prometheus/prometheus
[quay]: https://quay.io/repository/prometheus/prometheus
