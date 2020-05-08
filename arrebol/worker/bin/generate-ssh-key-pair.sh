#!/bin/sh
conf_file_path="$GOPATH/src/github.com/ufcg-lsd/arrebol-pb/arrebol/worker/worker-conf.json"
apt-get install jq

worker_id=`jq .id $conf_file_path`

key_name=$worker_id
key_size=${2-2048}
priv_key_path=$key_name.priv
openssl genrsa -out $key_name $key_size
openssl pkcs8 -topk8 -in $key_name -out $priv_key_path -nocrypt
openssl rsa -in $priv_key_path -outform PEM -pubout -out $key_name.pub
chmod 600 $priv_key_path
rm $key_name
