#!/bin/bash

set -eo pipefail

sanitize "${INPUT_USERNAME}" "username"
sanitize "${INPUT_PASSWORD}" "password"

echo "${INPUT_PASSWORD}" | docker login -u "${INPUT_USERNAME}" --password-stdin

function sanitize() {
    if [ -z "${1}" ]; then
        >&2 echo "Unable to find the ${2}. Did you set with.${2}?"
        exit 1
    fi
}
