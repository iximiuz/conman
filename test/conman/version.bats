#!/usr/bin/env bats

load helpers

function setup() {
    setup_test
    conmand_start
}

function teardown() {
    conmand_stop
}

@test "conmand version" {
    run conmanctl version
    debug $output
}

