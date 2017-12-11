#!/usr/bin/env bats

@test "environment create test" {
    l0 environment create test
}

@test "environment delete test" {
    l0 environment delete test
}

@test "job: list" {
    l0 job list
}

@test "job get (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job get $result
}

@test "job delete (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job delete $result
}
