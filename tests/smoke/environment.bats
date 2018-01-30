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

@test "environment create --scale 2 test3" {
    l0 environment create --scale 2 test3
}

@test "environment create --os windows --scale 1 test4" {
    l0 environment create --os windows --scale 1 test4
}

@test "environment link --bi-directional test1 test2" {
    l0 environment link --bi-directional test1 test2
}

@test "environment unlink --bi-directional test1 test2" {
    l0 environment unlink --bi-directional test1 test2
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

@test "environment scale test3 2" {
    l0 environment scale test3 2
}

@test "environment delete test1" {
    l0 --no-wait environment delete test1
}

@test "environment delete test2" {
    l0 --no-wait environment delete test2
}

@test "environment delete test3" {
    l0 --no-wait environment delete test3
}

@test "environment delete test4" {
    l0 environment delete test4
}
