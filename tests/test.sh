#!/usr/bin/env bash
set -e

for test in */ ; do
    echo
    echo "Testing ${test::-1}"
    echo "=========================="
    bats "$test/test.bats" || :
done
