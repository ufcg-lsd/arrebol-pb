#!/usr/bin/env bash

# Script useful for manage stack services life-cycle.
# See usage for commands & options.

set -o errexit;

readonly STACK_NAME=arrebol
readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

usage () {
    printf "usage: %s deploy | rm | [-h]" "$0"
}

define_params() {
    case $1 in
        deploy) shift
            docker stack deploy -c "${DIR}/"docker-stack.yml "${STACK_NAME}"
            ;;
        rm) shift
            docker stack rm "${STACK_NAME}"
            ;;
        -h | --help | *)
            usage;
            exit 0;
            ;;
    esac
}

define_params "$@"