#!/usr/bin/env bash

# Script useful for manage stack services life-cycle.
# See usage for commands & options.

set -o errexit;

readonly STACK_NAME=arrebol
readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
readonly VERSION=0.0.4

usage () {
    printf "usage: %s deploy | rm | [-h]" "$0"
}

define_params() {
    case $1 in
        build) shift
            docker build -t arrebol:"${VERSION}" -f stack/Dockerfile .
            ;;
        deploy) shift
            docker stack deploy -c "${DIR}/"docker-stack.yml "${STACK_NAME}"
            ;;
        clean) shift
            docker system prune -f
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