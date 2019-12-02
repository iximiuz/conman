#!/usr/bin/env bash
set -e

TEST_USERNS=${TEST_USERNS:-}

cd "$(dirname "$(readlink -f "$BASH_SOURCE")")"

bats $(pwd)

