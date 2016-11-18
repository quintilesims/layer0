setup () {
    cd "$BATS_TEST_DIRNAME"/../../cli/ || exit
    make build
    cd - || exit
    cp "$BATS_TEST_DIRNAME"/../../cli/l0 "$BATS_TEST_DIRNAME"/../common
    alias l0='../common/l0'
}
