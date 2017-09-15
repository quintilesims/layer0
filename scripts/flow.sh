#!/usr/bin/env bash
set -e

# kill all subprocesses on exit
trap 'if [ ! -z "$(jobs -pr)" ]; then kill $(jobs -pr); fi' EXIT

BULLET='\u2713'
GIT_HASH=$(git describe --tags)
LAYER0_PATH=$GOPATH/src/github.com/quintilesims/layer0

update_api() {
    echo "Updating API"
    pushd $LAYER0_PATH/api
        make release
    popd

    pushd $LAYER0_PATH/setup
        go run main.go set "$LAYER0_PREFIX" --input version="$GIT_HASH"
    popd
}

apply() {
    echo "Applying Changes"

    pushd $LAYER0_PATH/setup
        go run main.go apply "$LAYER0_PREFIX"
    popd
}

delete() {
    echo "Deleting Environments"
    environment_ids=$(l0 -o json environment list | jq -r .[].environment_id)
    for id in $environment_ids; do
        if [ "$id" != "api" ]; then
            echo -e $BULLET "$id"
            l0 environment delete $id --wait > /dev/null &
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

    echo "Deleting Jobs"
    job_ids=$(l0 -o json job list | jq -r .[].job_id)
    for id in $job_ids; do
        l0 job delete $id
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
            LAYER0_PREFIX=$OPTARG
            ;;
        *)
            exit 1
            ;;
    esac
done

if [ -z $LAYER0_PREFIX ]; then
    echo "LAYER0_PREFIX not set!"
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

