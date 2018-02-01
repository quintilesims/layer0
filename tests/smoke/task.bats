#!/usr/bin/env bats

@test "create" {
  l0 environment create env_name
  l0 deploy create ./common/Task.Dockerrun.aws.json dpl_name
  l0 task create env_name tsk_name1 dpl_name
  l0 task create --env c1:COMMAND="sleep 1" --env c2:COMMAND="sleep 2" env_name tsk_name2 dpl_name
}

@test "get" {
  l0 task get tsk_name1
  l0 task get tsk_name2
  l0 task get tsk_name*
}

@test "list" {
  l0 task list
}

@test "logs" {
  l0 task logs tsk_name1
  l0 task logs --tail 100 --start '2001-01-01 01:01' --end '2012-12-12 12:12' tsk_name1
}

@test "delete" {
  l0 task delete tsk_name1
  l0 task delete tsk_name2
  l0 deploy delete dpl_name
  l0 environment delete env_name
}

@test "deploy delete alpine:latest" {
    l0 deploy delete alpine:latest
}

# this deletes the remaining service(s), load balancer(s), and task(s)
@test "environment delete test" {
    l0 environment delete test
}
