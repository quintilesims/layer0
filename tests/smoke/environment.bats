#!/usr/bin/env bats

@test "create" {
  l0 environment create env_name1
  l0 environment create --user-data common/user_data.sh --os windows --type t2.small env_name2
}

@test "get" {
  l0 environment get env_name1
  l0 environment get env_name2
  l0 environment get env_name*
}

@test "list" {
  l0 environment list
}

@test "link" {
  l0 environment link --bi-directional env_name1 env_name2
}
 
@test "unlink" {
  l0 environment unlink --bi-directional env_name1 env_name2
}

@test "environment create --user-data common/user_data.sh test2" {
    l0 environment create --user-data common/user_data.sh test2 
}

@test "environment create --min-scale 2 --max-scale 5 test3" {
    l0 environment create --min-scale 2 --max-scale 5  test3
}

@test "environment create --os windows test4" {
    l0 environment create --os windows test4
}

@test "environment create test5" {
    l0 environment create test5
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

@test "environment set-scale --min-scale 1 --max-scale 5 test3" {
    l0 environment set-scale --min-scale 1 --max-scale 5 test3
}

@test "environment delete test1" {
    l0 --no-wait environment delete test1
}

@test "environment delete test2" {
    l0 --no-wait environment delete test2
}

@test "environment delete test3" {
    l0 --no-wait environment delete test3

@test "delete" {
  l0 environment delete env_name1
  l0 environment delete env_name2
}

@test "environment delete test5" {
    l0 environment delete -r test5
}

