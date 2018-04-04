#!/usr/bin/env bats

@test "create" {
  l0 environment create env_name
  l0 loadbalancer create env_name lb_name1
  l0 loadbalancer create --port 80:80/http --port 81:81/http --private env_name lb_name2
}

@test "get" {
  l0 loadbalancer get lb_name1
  l0 loadbalancer get lb_name2
  l0 loadbalancer get lb_name*
}

@test "list" {
  l0 loadbalancer list
}

@test "addport" {
  l0 loadbalancer addport lb_name1 8000:8000/http
}

@test "dropport" {
  l0 loadbalancer dropport lb_name1 8000
}

@test "healthcheck" {
  l0 loadbalancer healthcheck lb_name1
}

@test "delete" {
  l0 loadbalancer delete lb_name1
  l0 loadbalancer delete lb_name2
  l0 environment delete env_name
}
