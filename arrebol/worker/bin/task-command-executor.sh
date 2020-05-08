#!/bin/bash
set -e

execute_command() {
  command=$1
  /bin/bash -c "$command"
  exit $?
}

while getopts ":c:" opt; do
  case ${opt} in
    c)
      execute_command $OPTARG
      exit $?
      ;;
    \?)
      echo "Invalid option: $OPTARG" 1>&2
      ;;
    :)
      echo "Invalid option: $OPTARG requires an argument" 1>&2
      ;;
  esac
done