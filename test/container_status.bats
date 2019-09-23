#!/usr/bin/env bats

load helpers

function setup() {
    setup_test
    conmand_start
}

function teardown() {
    conmand_stop
}

@test "container status" {
    local cont_name="cont1"
    run conmanctl container create \
        --image "${TEST_ROOT}/data/rootfs_alpine/" \
        "${cont_name}" -- /bin/sleep 999
    [ $status -eq 0 ]

    local cont_id=$(echo $output | jq -r '.containerId')

    run conmanctl container status "${cont_id}"
    [ $status -eq 0 ]
    [ "${cont_id}" = $(echo $output | jq -r '.status.containerId') ]
    [ "CREATED" = $(echo $output | jq -r '.status.state') ]

    # TODO: kill all runc spawned by conmand
    run conmanctl container stop "${cont_id}"
    [ $status -eq 0 ]
}

