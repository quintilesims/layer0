#!/usr/bin/env bats

@test "environment create test" {
    l0 environment create test
}

@test "loadbalancer create --port 80:80/http test loadbalancer1" {
    l0 loadbalancer create --port 80:80/http test loadbalancer1
}

@test "deploy create guestbook" {
    l0 deploy create ./common/Service.Dockerrun.aws.json guestbook
}

@test "service create --loadbalancer loadbalancer1 test service1 guestbook:latest" {
    l0 service create --loadbalancer loadbalancer1 test service1 guestbook:latest
}

@test "service create --wait test service2 guestbook:latest" {
    l0 service create --wait test service2 guestbook:latest
}

@test "service create test service3 guestbook:latest" {
    l0 service create test service3 guestbook:latest
}

@test "service list" {
    l0 service list
}

@test "service get service1" {
    l0 service get service1
}

@test "service scale service1 2" {
    l0 service scale service1 2
}

@test "service scale --wait service2 2" {
    l0 service scale --wait service2 2
}

@test "deploy create guestbook" {
    l0 deploy create ./common/Service.Dockerrun.aws.json guestbook
}

@test "service update service1 guestbook:latest" {
    l0 service update service1 guestbook:latest
}

@test "service update --wait service2 guestbook:latest" {
    l0 service update --wait service2 guestbook:latest
}

@test "service update service3 guestbook:latest" {
    l0 service update service3 guestbook:latest
}

@test "service logs service1" {
    l0 service logs service1
}

@test "service logs --tail 100 service1" {
    l0 service logs --tail 100 service1
}

@test "service logs --start 01/01 01:01 --end 12/12 12:12 service1" {
    l0 service logs --start '01/01 01:01' --end '12/12 12:12' service1
}

# twice since we created 2 deploys named guestbook
@test "service: delete guestbook deploy" {
    l0 deploy delete guestbook:latest
    l0 deploy delete guestbook:latest
}

@test "service delete --wait service1" {
    l0 service delete --wait service1
}

# this deletes the remaining service(s) and loadbalancer(s)
@test "environment delete --wait test" {
    l0 environment delete --wait test
}
