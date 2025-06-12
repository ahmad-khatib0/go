#!/bin/bash

# build
go build .

# build for variety of OSs
GOOS="linux" go build

# list dependencies
go list -m all

# list available versions of the mux library
go list -m -versions github.com/gorilla/mux

# init the package manager
go mod init project-name

# verify the health of packages
go mod verify

# recompile packages
go mod tidy

# list which modules are depending on mux library
go mod why github.com/gorilla/mux

# list which dependencies are depending on each other
go mod graph

# change Go version in the mod file
go mod edit -go 1.7

# get dependencies into a vendor folder (like node_modules)
go mod vendor

# run from the vendor folder directly
go run -mod=vendor main.go

# check where you have race conditions
go run --race .

# see what versions of a module are available
go list -m -versions github.com/ahmad-khatib0/go/idiomatic-approach-book/simpletax

# output test coverage to a file
go test -v -cover -coverprofile=c.out

# run all benchmarks in a directory
go test -bench=.

# run a specific benchmark with CPU profiling
go test -bench BenchmarkGetIndex -cpuprofile cpu-books.out ./chapter08/performance

# run tests with race checker
go test -race

# run the fuzz test for 5 seconds
go test -fuzz FuzzGetSortedValues_ASC -fuzztime 5s ./chapter10/fragile-revised -v

# run specific tests
go test -run TestDivide ./ch04_test_suites/table -v
go test -run "^TestAdd" ./chapter02/calculator -cover -v
go test -run TestIndexIntegration ./chapter05/handlers -v -short

# coverage profile for a package
go test ./chapter02/calculator -coverprofile=calcCover.out

# run integration tests
go test -v -run Integration ./...

# using env var to target Integration tests (Go test does not have -long flag)
LONG=true go test -run TestIndexIntegration ./chapter05/handlers -v

# run BDD tests using Ginkgo
ginkgo -v ./chapter05/handlers # or using ./... operator

# View the file using the pprof command tool
go tool pprof performance.test cpu-books.out

# open the generated coverage file in HTML
go tool cover -html=calcCover.out

# Generate Ginkgo test suite
ginkgo bootstrap

# fix the error: go get -> module found but does not contain package
go clean -cache
go clean -modcache

# optionally force Go to fetch directly without proxy
export GOPROXY=direct
