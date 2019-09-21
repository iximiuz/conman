#!/usr/bin/env bats

load helpers

function setup() {
    setup_test
}

@test "conmand version" {
    conmand_start

    run $CONMCTL_BINARY version
    echo $output

    conmand_stop
}

