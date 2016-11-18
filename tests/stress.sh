#!/usr/bin/env bats
load ../common/common

@test "environment create test1" {
    l0 environment create test1
}

#!/bin/bash

for ((i=1;i<=100;i++));
do
   # your-unix-command-here
   echo $i
done
