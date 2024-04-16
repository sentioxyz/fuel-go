#!/bin/sh

cd $(dirname $0)

go run gen.go && gofmt -w ../
