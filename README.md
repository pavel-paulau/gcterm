# gcterm

[![codebeat badge](https://codebeat.co/badges/3ae2d27e-fc31-4189-8dc2-995589eb0e70)](https://codebeat.co/projects/github-com-pavel-paulau-gcterm)
[![Go Report Card](https://goreportcard.com/badge/github.com/pavel-paulau/gcterm)](https://goreportcard.com/report/github.com/pavel-paulau/gcterm)
[![Build Status](https://travis-ci.org/pavel-paulau/gcterm.svg?branch=master)](https://travis-ci.org/pavel-paulau/gcterm)
[![Coverage Status](https://coveralls.io/repos/github/pavel-paulau/gcterm/badge.svg?branch=master)](https://coveralls.io/github/pavel-paulau/gcterm?branch=master)

![Demo](http://imgur.com/hoh4b2O.png)

# Usage

First, enable [tracing of GC events](https://golang.org/pkg/runtime/) in your application using the GODEBUG variable:

    GODEBUG=gctrace=1

Redirect standard error to a file:

    $ myapp 2> stderr.log

Now you should be able to visualize GC events in your terminal:

    $ tail -f stderr.log | gcterm

Currently, only applications built with Go 1.6 and Go 1.7 are supported.
