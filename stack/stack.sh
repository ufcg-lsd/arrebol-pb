#!/usr/bin/env bash

# Script useful for manage stack services life-cycle.
# See usage for commands & options.

set -o errexit;

readonly STACK_NAME=arrebol
readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
readonly VERSION="$(cat < .env | grep VERSION_TAG= | cut -d= -f2)"
readonly DOCKERHUB_REPO=emanueljoivo/arrebol

usage ()
{
    printf "usage: %s deploy | rm | [-h]" "$0"
}

arrebol_build()
{
  local DOCKERFILE_DIR=stack/Dockerfile

  docker build -t "${DOCKERHUB_REPO}":"${VERSION}" \
            --file "${DOCKERFILE_DIR}" .
}

stack_deploy()
{
  docker stack deploy -c "${DIR}/"docker-stack.yml "${STACK_NAME}"
}

stack_publish()
{
  docker push "${DOCKERHUB_REPO}":"${VERSION}"
}

stack_clean()
{
  docker system prune -f
}

stack_rm()
{
  docker stack rm "${STACK_NAME}"
}

release() {
  local VERSION_TAG=$1
  local VERSION_NAME=$2
  sed -i "/^VERSION_TAG/c\VERSION_TAG=${VERSION_TAG}" ./.env
  sed -i "/^VERSION_NAME/c\VERSION_NAME=${VERSION_NAME}" ./.env
  arrebol_build
  stack_publish
}

define_params()
{
    case $1 in
        release) shift
            release "$@"
            ;;
        build) shift
            arrebol_build
            ;;
        deploy) shift
            stack_deploy
            ;;
        publish) shift
            arrebol_build
            stack_publish
            ;;
        clean) shift
            stack_clean
            ;;
        rm) shift
            stack_rm
            ;;
        -h | --help | *)
            usage;
            exit 0;
            ;;
    esac
}

define_params "$@"