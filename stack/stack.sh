#!/usr/bin/env bash

# Script useful for manage stack services life-cycle.
# See usage for commands & options.

set -o errexit;

readonly STACK_NAME=arrebol
readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
readonly VERSION=0.0.4
readonly DOCKERHUB_REPO=emanueljoivo/arrebol

usage () {
    printf "usage: %s deploy | rm | [-h]" "$0"
}

build_arrebol() {
  local DOCKERFILE_DIR=stack/Dockerfile

  docker build -t "${DOCKERHUB_REPO}":"${VERSION}" \
            --file "${DOCKERFILE_DIR}" .
}

define_params() {
    case $1 in
        build) shift
            build_arrebol
            ;;
        deploy) shift
            docker stack deploy -c "${DIR}/"docker-stack.yml "${STACK_NAME}"
            ;;
        publish) shift
            docker push "${DOCKERHUB_REPO}":"${VERSION}"
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