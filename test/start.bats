#!/usr/bin/env bats

load helpers

function setup() {
    setup_test
}

@test "conmand start" {
    conmand_start
    conmand_stop
}

