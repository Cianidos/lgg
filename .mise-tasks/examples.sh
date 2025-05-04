#!/usr/bin/env bash
#MISE description="Build the example and run"

run_example() {
    local example_name=$1
    echo ""
    echo "------ COMPILE example $example_name"
    go run cmd/main.go -- examples/"$example_name".lgo >examples/"$example_name".go
    echo "------ RUN $example_name"
    go run examples/"$example_name".go
    echo "------ END $example_name end"
}

run_example hello_world
run_example hello_world2
run_example empty_main
run_example defun_and_call
