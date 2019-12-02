#!/usr/bin/env bats

load helpers

function setup() {
    setup_test
    conmand_start
}

function teardown() {
    conmand_stop
}

@test "conmand restore" {
    # Create container 1
    run conmanctl container create \
        --image "${TEST_ROOT}/data/rootfs_alpine/" \
        cont1 -- /bin/sleep 100
    [ $status -eq 0 ]

    local cont_id1=$(jq -r '.containerId' <<< $output)

    # Create, then start container 2
    run conmanctl container create \
        --image "${TEST_ROOT}/data/rootfs_alpine/" \
        cont2 -- /bin/sleep 200
    [ $status -eq 0 ]

    local cont_id2=$(jq -r '.containerId' <<< $output)
    run conmanctl container start "${cont_id2}"
    [ $status -eq 0 ]

    # Create, start, then stop container 3
    run conmanctl container create \
        --image "${TEST_ROOT}/data/rootfs_alpine/" \
        cont3 -- /bin/sleep 300
    [ $status -eq 0 ]

    local cont_id3=$(jq -r '.containerId' <<< $output)
    run conmanctl container start "${cont_id3}"
    [ $status -eq 0 ]

    run conmanctl container stop "${cont_id3}"
    [ $status -eq 0 ]

    # Assert list
    run conmanctl container list
    [ $status -eq 0 ]
    
    local list1=$(jq -rc . <<< $output)
    debug $list1

    local state1=$(jq -r --arg ID "$cont_id1" '.containers[] | select(.id == $ID).state' <<< $list1)
    local state2=$(jq -r --arg ID "$cont_id2" '.containers[] | select(.id == $ID).state' <<< $list1)
    local state3=$(jq -r --arg ID "$cont_id3" '.containers[] | select(.id == $ID).state' <<< $list1)

    [ "CREATED" = "${state1}" ]
    [ "RUNNING" = "${state2}" ]
    [ "EXITED" = "${state3}" ]

    # restart conmand
    conmand_restart

    # assert same list
    run conmanctl container list
    [ $status -eq 0 ]

    local list2=$(jq -rc . <<< $output)
    debug $list2

    [ $list1 = $list2 ]

    # TODO: kill all runc spawned by conmand
    run conmanctl container stop "${cont_id1}"
    run conmanctl container stop "${cont_id2}"
}
