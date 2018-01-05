#!/usr/bin/env bats

@test "environment create test" {
    l0 environment create test
}

@test "deploy create alpine" {
    l0 deploy create ./common/Task.Dockerrun.aws.json alpine
}

@test "task create --wait test task1 alpine:latest" {
    l0 task create --wait test task1 alpine:latest
}

@test "task create --env alpine:key=val --copies 3 test task2 alpine:latest" {
    l0 task create --env alpine:key=val --copies 3 test task2 alpine:latest
}

@test "task list" {
    l0 task list
}

@test "task list --all" {
    l0 task list --all
}

@test "task get task1" {
    l0 task get task1
}

@test "task get t*" {
    l0 task get t\*
}

@test "task logs task1" {
    l0 task logs task1
}

@test "task logs --tail 100 task1" {
    l0 task logs --tail 100 task1
}

@test "task logs --start 2001-01-01 01:01 --end 2012-12-12 12:12 task1" {
    l0 task logs --start '2001-01-01 01:01' --end '2012-12-12 12:12' task1
}

@test "task delete task1" {
    l0 task delete task1
}

@test "deploy delete alpine:latest" {
    l0 deploy delete alpine:latest
}

# this deletes the remaining service(s), loadbalancer(s), and task(s)
@test "environment delete --wait test" {
    l0 environment delete --wait test
}
