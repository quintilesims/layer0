#!/usr/bin/env bats

@test "admin sql" {
    l0 admin sql
}

@test "admin version" {
    l0 admin version
}

@test "admin debug" {
    l0 admin debug
}
