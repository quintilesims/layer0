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
    l0 job get $(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')
}

@test "job: logs (most recent)" {
    l0 job logs $(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')
}

@test "job: logs --tail 100 (most recent)" {
    l0 job logs --tail 100 $(l0 -o json job list | jq -r 'max_by(.time_created) | .job_id')
}
