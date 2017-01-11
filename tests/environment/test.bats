#!/usr/bin/env bats

@test "environment create test1" {
    l0 environment create test1
}

@test "environment list" {
    l0 environment list
}

@test "environment get test1" {
    l0 environment get test1
}

@test "environment get t" {
    l0 environment get t\*
}

@test "environment create --user-data environment/user_data.sh test2" {
    l0 environment create --user-data environment/user_data.sh test2 
}

@test "environment create --min-count 2 test3" {
    l0 environment create --min-count 2 test3
}

@test "environment list" {
    l0 environment list
}

@test "environment get test2" {
    l0 environment get test2
}

@test "environment setmincount test3 0" {
    l0 environment setmincount test3 0
}

@test "environment setmincount test3 3" {
    l0 environment setmincount test3 3
}

@test "environment delete test1" {
    l0 environment delete test1
}

@test "environment delete test2" {
    l0 environment delete test2
}

@test "environment delete --wait test3" {
    l0 environment delete --wait test3
}
