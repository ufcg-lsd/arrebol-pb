#!/bin/bash

check_docker_installation() {
  sudo docker --version
  if [[ "$?" != 0 ]]; then
    sudo apt-get install docker.io
    sudo systemctl start docker
    sudo systemctl enable docker
  fi
}

check_docker_installation

TAG=${1-latest}
IMAGE=ufcg-lsd/arrebol:$TAG
ID=worker_id=`jq .id $conf_file_path`
CONTAINER_NAME=arrebol-worker-$ID
PROJECT_PATH=/go/src/github.com/ufcg-lsd/arrebol-pb/arrebol
sudo docker pull $IMAGE
sudo docker stop $CONTAINER_NAME
sudo docker rm $CONTAINER_NAME
sudo docker run --name $CONTAINER_NAME -tdi $IMAGE


