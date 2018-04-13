#!/usr/bin/env bat

load common/test_helpers

@test "create" {
  l0 environment create --scale 1 env_name
  waitForEnvironmentToScale env_name
  l0 deploy create ./common/Task_stateful.Dockerrun.aws.json dpl_name_stateful
  l0 deploy create ./common/Task_stateless.Dockerrun.aws.json dpl_name_stateless
  l0 task create env_name tsk_name_stateless dpl_name_stateless:latest
  l0 task create --stateful env_name tsk_name_stateful1 dpl_name_stateful:latest
  l0 task create --env c1:COMMAND="sleep 1" --env c2:COMMAND="sleep 2" --stateful env_name tsk_name_stateful2 dpl_name_stateful:latest
}

@test "get" {
  l0 task get tsk_name_stateless
  l0 task get tsk_name_stateful1
  l0 task get tsk_name*
}

@test "list" {
  l0 task list
}

@test "logs" {
  l0 task logs tsk_name_stateless
  l0 task logs --tail 100 --start '2001-01-01 01:01' --end '2012-12-12 12:12' tsk_name_stateful1
}

@test "delete" {
  l0 task delete tsk_name_stateless
  l0 task delete tsk_name_stateful1
  l0 task delete tsk_name_stateful2
  l0 deploy delete dpl_name_stateful:latest
  l0 deploy delete dpl_name_stateless:latest
  l0 environment delete env_name
}