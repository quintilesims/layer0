#!/usr/bin/env bats
load ../common/common

setup(){
    create_cert
}

teardown(){
    delete_cert
}

@test "certificate create certificate1" {
    l0 certificate create certificate1 www.example.com.cert www.example.com.key
}

@test "certificate list" {
    l0 certificate list
}

@test "certificate get certificate1" {
    l0 certificate get certificate1
}

@test "certificate delete certificate1" {
    l0 certificate delete certificate1
}
