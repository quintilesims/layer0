#!/usr/bin/env bats
load ../common/common

@test "admin sql" {
    l0 admin sql
}

@test "admin version" {
    l0 admin version
}

@test "admin debug" {
    l0 admin debug
}
