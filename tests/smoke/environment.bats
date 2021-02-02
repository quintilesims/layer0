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

@test "environment create --user-data common/user_data.sh test2" {
    l0 environment create --user-data common/user_data.sh test2 
}

@test "environment create --min-count 2 --max-count 2 --target-cap-size 100 test3" {
    l0 environment create --min-count 2 --max-count 2 --target-cap-size 100 test3
}

@test "environment create --os windows test4" {
    l0 environment create --os windows test4
}

@test "environment link test3 test4" {
    l0 environment link test3 test4
}

@test "environment unlink test3 test4" {
    l0 environment unlink test3 test4
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
    l0 environment delete test3
}

@test "environment delete --wait test4" {
    l0 environment delete --wait test4
}
