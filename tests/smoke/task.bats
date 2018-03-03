#!/usr/bin/env bat

@test "create" {
  l0 environment create --scale 1 env_name
  l0 deploy create ./common/Task_stateful.Dockerrun.aws.json dpl_name_stateful
  l0 deploy create ./common/Task_stateless.Dockerrun.aws.json dpl_name_stateless
  l0 task create --stateful env_name tsk_name_stateful dpl_name_stateful
  l0 task create env_name tsk_name_stateless dpl_name_stateless
  l0 task create --env c1:COMMAND="sleep 1" --env c2:COMMAND="sleep 2" env_name tsk_name_stateful2 dpl_name_stateful
}

@test "get" {
  l0 task get tsk_name_stateful
  l0 task get tsk_name_stateless
  l0 task get tsk_name_stateless2
  l0 task get tsk_name*
}

@test "list" {
  l0 task list
}

@test "logs" {
  l0 task logs tsk_name_stateful
  l0 task logs --tail 100 --start '2001-01-01 01:01' --end '2012-12-12 12:12' tsk_name_stateful
}

@test "delete" {
  l0 task delete tsk_name_stateful
  l0 task delete tsk_name_stateful2
  l0 task delete tsk_name_stateless
  l0 deploy delete dpl_name_stateful
  l0 deploy delete dpl_name_stateless
  l0 environment delete env_name
}

@test "deploy delete alpine:latest" {
    l0 deploy delete alpine:latest
}

# this deletes the remaining service(s), load balancer(s), and task(s)
@test "environment delete test" {
    l0 environment delete test
}
