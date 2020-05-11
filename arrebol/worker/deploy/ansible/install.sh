#!/bin/bash

check_ansible_installation() {
  ansible --version
  if [[ "$?" != 0 ]]; then
    sudo apt-get install ansible
  fi
}

check_ansible_installation

DEPLOY_CONF_FILE="$(pwd)/deploy.conf"
ANSIBLE_HOSTS_FILE="$(pwd)/hosts"
ANSIBLE_CFG_FILE="$(pwd)/ansible.cfg"

# setup hosts
WORKER_HOST_IP=$(grep ^host DEPLOY_CONF_FILE | awk -F "=" '{print $2}')
PATTERN_HELPER="\[worker-machine\]"
WORKER_HOST_IP_PATTERN=$(grep -A1 $PATTERN_HELPER $ANSIBLE_HOSTS_FILE | tail -n 1)
sed -i "s/$WORKER_HOST_IP_PATTERN/$WORKER_HOST_IP/" $ANSIBLE_HOSTS_FILE

ANSIBLE_SSH_PRIVATE_KEY_FILE=$(grep ^private_key_path $DEPLOY_CONF_FILE | awk -F "=" '{print $2}')
PRIVATE_KEY_FILE_PATH_PATTERN="ansible_ssh_private_key_file"
sed -i "s#.*$PRIVATE_KEY_FILE_PATH_PATTERN=.*#$PRIVATE_KEY_FILE_PATH_PATTERN=$ANSIBLE_SSH_PRIVATE_KEY_FILE#g" $ANSIBLE_HOSTS_FILE


#setup ansible.cfg
REMOTE_USER=$(grep ^remote_host_user $DEPLOY_CONF_FILE | awk -F "=" '{print $2}')
PATTERN_HELPER="remote_user"
sed -i "s#.*$PATTERN_HELPER = .*#$PATTERN_HELPER = $REMOTE_USER#g" $ANSIBLE_CFG_FILE

ansible-playbook -vvv ./deploy.yml