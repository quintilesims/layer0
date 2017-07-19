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

@test "loadbalancer create --healthcheck-target TCP:80 --healthcheck-interval 30 --healthcheck-timeout 5 --healthcheck-healthy-threshold 2 --healthcheck-unhealthy-threshold 2 loadbalancer3" {
    l0 loadbalancer create --healthcheck-target TCP:80 --healthcheck-interval 30 --healthcheck-timeout 5 --healthcheck-healthy-threshold 2 --healthcheck-unhealthy-threshold 2 test loadbalancer3
}

@test "loadbalancer healthcheck loadbalancer3" {
    l0 loadbalancer healthcheck loadbalancer3
}

@test "loadbalancer healthcheck --set-target TCP:88 --set-interval 45 --set-timeout 10 --set-healthy-threshold 5 --set-unhealthy-threshold 3 loadbalancer3" {
    l0 loadbalancer healthcheck --set-target TCP:88 --set-interval 45 --set-timeout 10 --set-healthy-threshold 5 --set-unhealthy-threshold 3 loadbalancer3
}

@test "loadbalancer delete --wait loadbalancer3" {
    l0 loadbalancer delete --wait loadbalancer3
}

@test "loadbalancer create --port 80:80/http test loadbalancer4" {
    l0 loadbalancer create --port 80:80/http test loadbalancer4
}

@test "deploy create guestbook" {
    l0 deploy create ./common/Service.Dockerrun.aws.json guestbook
}

@test "service create --loadbalancer loadbalancer4 test service1 guestbook:latest" {
    l0 service create --loadbalancer loadbalancer4 test service1 guestbook:latest
}

@test "loadbalancer list" {
    l0 loadbalancer list
}

@test "loadbalancer get loadbalancer4" {
    l0 loadbalancer get loadbalancer4
}

@test "deploy delete guestbook:latest" {
    l0 deploy delete guestbook:latest
}

# this deletes the remaining service(s) and loadbalancer(s)
@test "environment delete --wait test" {
    l0 environment delete --wait test
}
