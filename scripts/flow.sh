#!/usr/bin/env bash
set -e

# kill all subprocesses on exit
trap 'if [ ! -z "$(jobs -pr)" ]; then kill $(jobs -pr); fi' EXIT

BULLET='->'
GIT_HASH=$(git describe --tags)
LAYER0_PATH=$GOPATH/src/github.com/quintilesims/layer0

update_api() {
    echo "Updating API"
    pushd $LAYER0_PATH/api
        make release
    popd

    pushd $LAYER0_PATH/setup
        go run main.go set "$LAYER0_INSTANCE" --input version="$GIT_HASH"
    popd
}

apply() {
    echo "Applying Changes"

    pushd $LAYER0_PATH/setup
        go run main.go apply "$LAYER0_INSTANCE"
    popd
}

delete() {
    echo "Deleting Load Balancers"
    load_balancer_ids=$(l0 -o json loadbalancer list | jq -r .[].load_balancer_id)
    for id in $load_balancer_ids; do
        if [ "$id" != "api" ]; then
            echo -e $BULLET "$id"
            l0 loadbalancer delete $id > /dev/null &
        fi
    done

    echo "Deleting Services"
    service_ids=$(l0 -o json service list | jq -r .[].service_id)
    for id in $service_ids; do
        if [ "$id" != "api" ]; then
            echo -e $BULLET "$id"
            l0 service delete $id > /dev/null &
        fi
    done

    echo "Deleting Environments"
    environment_ids=$(l0 -o json environment list | jq -r .[].environment_id)
    for id in $environment_ids; do
        if [ "$id" != "api" ]; then
            echo -e $BULLET "$id"
            l0 environment delete $id > /dev/null &
        fi
    done

    wait

    echo "Deleting Deploys"
    deploy_ids=$(l0 -o json deploy list --all | jq .[] | jq 'select(.deploy_name != "")' | jq -r .deploy_id)
    for id in $deploy_ids; do
        echo -e $BULLET "$id"
        l0 deploy delete $id > /dev/null
        echo -e $BULLET "$id"
    done

}

usage() {
    echo "Usage: flow [OPTIONS...] ARGUMENTS...
Build and push Docker images from your current Layer0 code.
Update your Layer0 configuration to run the new images.

Options:
    -h, --help      Show help menu
    -p, --prefix    Specify which Layer0 prefix to apply changes on

Arguments:
    api             Build and update docker image for the Layer0 API
    delete          Delete all entities in a Layer0
"
}

while getopts "hp:" option; do
    case "$option" in
        h)
            usage
            exit 0
            ;;
        p)
            LAYER0_INSTANCE=$OPTARG
            ;;
        *)
            exit 1
            ;;
    esac
done

if [ -z $LAYER0_INSTANCE ]; then
    echo "LAYER0_INSTANCE not set!"
    exit 1
fi

for i in ${@:$OPTIND}
do
    case $i in
        help)
            usage
            exit
            ;;
        api)
            update_api
            should_apply=true
            ;;
        delete)
            delete
            ;;
        *)
            echo "Incorrect Usage. Unknown argument '"$i"'"
            exit 1
            ;;
    esac
done

if [ ! -z $should_apply ]; then
    apply
fi

