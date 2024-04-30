#!/bin/sh

cd $(dirname $0)

go run *.go && gofmt -w ../
