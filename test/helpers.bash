#!/usr/bin/env bash

TEST_ROOT=$(dirname "$(readlink -f "$BASH_SOURCE")")
PROJECT_ROOT=$(cd "$TEST_ROOT/.."; pwd -P)
BIN_ROOT="${PROJECT_ROOT}/bin";

RUNTIME_PATH="/usr/bin/runc"
RUNTIME_ROOT="/run/conman-runc"

CONMAND_BINARY="${BIN_ROOT}/conmand"
CONMAND_LOG="${TEST_ROOT}/conmand.log"
CONMAND_DIR=
CONMAND_PID=

CONMANCTL_BINARY="${BIN_ROOT}/conmanctl"

function conmand_start() {
    conmand_setup

    debug "starting conmand in ${CONMAND_DIR}"
    $CONMAND_BINARY \
        --lib-root "${CONMAND_DIR}/var/lib/conman" \
        --run-root "${CONMAND_DIR}/run/conman" \
        --runtime-path "${RUNTIME_PATH}" \
        --runtime-root "${CONMAND_DIR}/${RUNTIME_ROOT}" \
        --listen "${CONMAND_DIR}/run/conmand.sock" \
        &> >(tee $CONMAND_LOG) & CONMAND_PID=$!
    debug "conmand PID ${CONMAND_PID}"

    conmand_wait
}

function conmand_restart() {
    debug "starting conmand in ${CONMAND_DIR}"

    if [[ "$CONMAND_PID" -eq 0 ]] ; then
        return 1
    fi

    kill $CONMAND_PID && wait $CONMAND_PID || true
    $CONMAND_BINARY \
        --lib-root "${CONMAND_DIR}/var/lib/conman" \
        --run-root "${CONMAND_DIR}/run/conman" \
        --runtime-path "${RUNTIME_PATH}" \
        --runtime-root "${CONMAND_DIR}/${RUNTIME_ROOT}" \
        --listen "${CONMAND_DIR}/run/conmand.sock" \
        &> >(tee $CONMAND_LOG) & CONMAND_PID=$!
    debug "conmand PID ${CONMAND_PID}"

    conmand_wait
}

function conmand_setup() {
    CONMAND_DIR=$(mktemp --directory --tmpdir="/tmp" conman-test-run.XXXXXX)
}

function conmand_wait() {
    retry 5 1 conmanctl version
}

function conmand_stop() {
    if [[ "$CONMAND_PID" -eq 0 ]] ; then
        return 1
    fi

    debug "terminating conmand ${CONMAND_PID}"

    kill $CONMAND_PID && wait $CONMAND_PID || true
    CONMAND_PID=
    CONMAND_DIR=
}

function conmanctl() {
    "${CONMANCTL_BINARY}" --host "${CONMAND_DIR}/run/conmand.sock" "$@"
}

function setup_test() {
    debug ' +++++++++++++++++++++++++++++++++++++++++++ '
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
    local lines=$@
    for l in "${lines[@]}"; do
        echo "# $l" >&3
    done
}

