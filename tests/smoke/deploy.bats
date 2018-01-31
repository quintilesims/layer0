#!/usr/bin/env bats

@test "create" {
  l0 deploy create ./common/Service.Dockerrun.aws.json dpl_name1
  l0 deploy create ./common/Service.Dockerrun.aws.json dpl_name2
}

@test "get" {
  l0 deploy get dpl_name1
  l0 deploy get dpl_name2
  l0 deploy get dpl_name*
}

@test "list" {
  l0 deploy list
  l0 deploy list --all
}

@test "delete" {
  l0 deploy delete dpl_name1
  l0 deploy delete dpl_name2
}
