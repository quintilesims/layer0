#!/usr/bin/env bash
set -e

dirs=(admin deploy environment job load_balancer service task)

for dir in "${dirs[@]}"; do
    echo
    echo "Testing $dir"
    echo "====================="
    bats "$dir/test.bats"
done
