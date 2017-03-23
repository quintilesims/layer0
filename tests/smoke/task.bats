#!/usr/bin/env bats

@test "environment create test" {
    l0 environment create test
}

@test "deploy create guestbook" {
    l0 deploy create ./common/Dockerrun.aws.json guestbook
}

@test "task create --wait test task1 guestbook" {
    l0 task create --wait test task1  guestbook
}

@test "task create --copies 2 --env guestbook:key=val test task2 guestbook" {
    l0 task create --copies 2 --env guestbook:key=val test task2 guestbook 
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

@test "task delete task1" {
    l0 task delete task1
}

@test "deploy delete guestbook" {
    l0 deploy delete guestbook
}

# this deletes the remaining service(s) and loadbalancer(s)
@test "environment delete --wait test" {
    l0 environment delete --wait test
}
