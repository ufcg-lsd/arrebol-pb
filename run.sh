#!/usr/bin/env bash

set -o errexit
set -o xtrace

_install_docker() {
    wget -qO- https://get.docker.com | sh

    /etc/init.d/docker status

    groupadd docker

    usermod -aG docker $USER

    service docker restart

    su ${USER}
}

main() {
  if [[ ${EUID} != 0 ]]; then
      echo "Must be executed as root"
  else
       #_install_docker


      exit 0
  fi
}

main