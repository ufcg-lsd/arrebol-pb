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
  local DB_ADDR=127.0.0.1
  local DB_USER=arrebol-admin
  local DB_PASSWD=postgres
  local DB_NAME=arrebol-db
  local DB_PORT=5432
  local WORKER_AMOUNT=5

  local DOCKERFILE_DIR=stack/Dockerfile

  docker build -t "${DOCKERHUB_REPO}":"${VERSION}" \
            --build-arg DB_ADDR="${DB_ADDR}" \
            --build-arg DB_USER="${DB_USER}" \
            --build-arg DB_PASSWD="${DB_PASSWD}" \
            --build-arg DB_NAME="${DB_NAME}" \
            --build-arg DB_PORT="${DB_PORT}" \
            --build-arg WORKER_AMOUNT="${WORKER_AMOUNT}" \
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