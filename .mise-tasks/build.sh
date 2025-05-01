#!/usr/bin/env bash
#MISE description="Build the example and run"
go run cmd/main.go -- examples/hello_world.lgo > examples/hello_world.go
go run examples/hello_world.go
