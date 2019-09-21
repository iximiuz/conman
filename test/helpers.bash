#!/usr/bin/env bash

TEST_ROOT=$(dirname "$(readlink -f "$BASH_SOURCE")")
PROJECT_ROOT=$(cd "$TEST_ROOT/.."; pwd -P)
BIN_ROOT="${PROJECT_ROOT}/bin";

CONMAND_BINARY="${BIN_ROOT}/conmand"
CONMAND_LOG="${TEST_ROOT}/conmand.log"
CONMAND_PID=0

CONMCTL_BINARY="${BIN_ROOT}/conmanctl"

function conmand_start() {
    debug "starting conmand"
    $CONMAND_BINARY &> >(tee $CONMAND_LOG) & CONMAND_PID=$!
    debug "conmand PID ${CONMAND_PID}"

    conmand_wait
}

function conmand_wait() {
    retry 10 1 $CONMCTL_BINARY version
}

function conmand_stop() {
    if [[ "$CONMAND_PID" -eq 0 ]] ; then
        return 1
    fi

    kill $CONMAND_PID 
    wait $CONMAND_PID || true
    CONMAND_PID=0
}

function setup_test() {
    debug ' +++++++++++++++++++++++++++++++++++ '
}

function retry() {
    local attempts=$1
    shift
    local delay=$1
    shift
    local i

    for ((i=0; i < attempts; i++)); do
        run "$@"
        if [[ "$status" -eq 0 ]] ; then
            return 0
        fi
        sleep $delay
    done

    return $status
}

function debug() {
    local lines=$1
    for i in "${lines[@]}"; do
        echo "# $i" >&3
    done
}

