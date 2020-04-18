#!/usr/bin/env bash

export PATH=$PATH:$(go env GOPATH)/bin
go get -u golang.org/x/lint/golint
golint -set_exit_status $(go list ./... | grep -v /vendor/)