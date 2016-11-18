check_service_running() {
    service=$1
    scale=$2
    current_scale=$(l0 -o json service get "$service" | jq '.[].running_count' -r)
    if [ "$current_scale" -eq "$scale" ]; then
        echo "OK"
    fi
}

wait_service_running(){
    service=$1
    scale=$2
    wait_for 300 check_service_running "$service" "$scale"
}

wait_for() {
    timeout=$1
    shift
    fcn=( "$@" )
    runtime=0
    sleep_duration=10
    while [ "$runtime" -lt "$timeout" ]; do
        out=$($fcn)
        if [[ $out =~ .*OK ]]; then # filter to the last line and check for OK
            exit 0
        fi
        runtime=($runtime+$sleep_duration)
        sleep $sleep_duration
    done
    exit 1
}

del_service() {
    service=$1
    out=$(l0 service delete "$service")
    # approx match the output because of the spinner
    if [ $? -eq 0 ] || [[ $out =~ "No service found with that Name or ID" ]]; then
        echo "OK"
    fi
}

del_lb() {
    lb=$1
    out=$(l0 loadbalancer delete "$lb")
    # approx match the output because of the spinner
    if [ $? -eq 0 ] || [[ $out =~ "No load_balancer found with that Name or ID" ]]; then
        echo "OK"
    fi
}

env_create() {
    env_name=$1
    output=$(l0 environment create "$env_name")
    if [[ $output = *"Environment $env_name exists"* ]]; then
        echo "OK"
    fi
}

env_delete() {
    env_name=$1
    output=$(l0 environment delete "$env_name")
    # approx match the output because of the spinner
    if [[ $output = *"No environment found with that Name or ID"* ]]; then
        echo "OK"
    fi
}

create_cert() {
    openssl req \
        -new \
        -newkey rsa:4096 \
        -days 365 \
        -nodes \
        -x509 \
        -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com" \
        -keyout www.example.com.key \
        -out www.example.com.cert
}

delete_cert() {
    rm www.example.com.key && rm www.example.com.cert
}

get_lb_url() {
    lb=$1
    runtime=0
    sleep_duration=10
    timeout=60
    while [ $runtime -lt $timeout ]; do
        maybe_url=$(l0 -o json loadbalancer get $lb | jq '.[].url' -r -M)
        if [[ $maybe_url != "null" ]]; then
            echo $maybe_url
            break
        fi
        runtime=$(($runtime+$sleep_duration))
        sleep $sleep_duration
    done
}
