#!/usr/bin/env bats

certificate_name="l0-$LAYER0_PREFIX-api"

@test "create environment test" {
    l0 environment create test
}

@test "loadbalancer create test loadbalancer1" {
    l0 loadbalancer create test loadbalancer1
}

@test "loadbalancer list" {
    l0 loadbalancer list
}

@test "loadbalancer get loadbalancer1" {
    l0 loadbalancer get loadbalancer1
}

@test "loadbalancer addport loadbalancer1 8000:8000/http" {
    l0 loadbalancer addport loadbalancer1 8000:8000/http
}

@test "loadbalancer dropport loadbalancer1 8000" {
    l0 loadbalancer dropport loadbalancer1 8000
}

@test "loadbalancer delete --wait loadbalancer1" {
    l0 loadbalancer delete --wait loadbalancer1
}

@test "loadbalancer create --port 80:80/http --port 443:443/https --private --certificate $certificate_name loadbalancer2" {
    l0 loadbalancer create --port 80:80/http --port 443:443/https --private --certificate $certificate_name test loadbalancer2
}

@test "loadbalancer list" {
    l0 loadbalancer list
}

@test "loadbalancer get loadbalancer2" {
    l0 loadbalancer get loadbalancer2
}

@test "loadbalancer get l\*" {
    l0 loadbalancer get l\*
}

@test "loadbalancer delete --wait loadbalancer2" {
    l0 loadbalancer delete --wait loadbalancer2
}

@test "loadbalancer create --port 80:80/http test loadbalancer3" {
    l0 loadbalancer create --port 80:80/http test loadbalancer3
}

@test "deploy create guestbook" {
    l0 deploy create ./deploy/Guestbook.dockerrun.aws.json guestbook
}

@test "service create --loadbalancer loadbalancer3 test service1 guestbook" {
    l0 service create --loadbalancer loadbalancer3 test service1 guestbook
}

@test "loadbalancer list" {
    l0 loadbalancer list
}

@test "loadbalancer get loadbalancer3" {
    l0 loadbalancer get loadbalancer3
}

@test "deploy delete guestbook" {
    l0 deploy delete guestbook
}

# this deletes the remaining service(s) and loadbalancer(s)
@test "environment delete --wait test" {
    l0 environment delete --wait test 
}

@test "certificate delete certificate1" {
    delete_cert
    l0 certificate delete certificate1
}
