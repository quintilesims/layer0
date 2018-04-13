#!/bin/bash

waitForEnvironmentToScale() {
    echo "waiting for environment to reach desired scale"

    # 300 seconds / sleep duration of 5 = 60
    for value in {1..60}
    do
        output=$(l0 -o json environment get env_name | jq '.[0]' | jq '.desired_size == .current_size')
        if [ "$output" = "true" ]; then
            echo "current has reached desired scale"
            break
        fi

        echo "waiting..."
        sleep 5
    done
}