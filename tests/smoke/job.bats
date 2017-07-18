#!/usr/bin/env bats

@test "environment create test" {
    l0 environment create test
}

@test "environment delete --wait test" {
    l0 environment delete -wait test
}

@test "job: list" {
    l0 job list
}

@test "job get (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job get $result
}

@test "job: logs (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job logs $result
}

@test "job: logs --tail 100 (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job logs $result
}

@test "job logs --start 2001-01-01 01:01 --end 2012-12-12 12:12 (most recent)" {
    result="$(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')"
    l0 job logs --start '2001-01-01 01:01' --end '2012-12-12 12:12' $result
}

