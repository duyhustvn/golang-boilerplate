#!/usr/bin/env bash

NUMBER_OF_VMS=3

echo "create rke2 config.yaml"

token_path="/var/lib/rancher/rke2/server/token"

token=$(vagrant ssh vm1 -c "sudo cat $token_path")

echo "token: $token"


for ((i = 2; i <= $NUMBER_OF_VMS; i++)); do

  vagrant ssh vm$i -c 'sudo mkdir -p /etc/rancher/rke2'

   node_name=node-master-0$i
   vagrant ssh vm$i -c "
cat <<EOF > ~/config.yaml
node-name: $node_name
token: $token
debug: true
# disable: rke2-ingress-nginx
cni:
  - canal
disable-cloud-controller: true
enable-servicelb: true
kube-apiserver-arg:
  - 'default-not-ready-toleration-seconds=30'
  - 'default-unreachable-toleration-seconds=30'
EOF
"
  vagrant ssh vm$i -c 'sudo mv config.yaml /etc/rancher/rke2'

  vagrant ssh vm1 -c 'sudo systemctl enable rke2-server'
  vagrant ssh vm$i -c 'sudo systemctl start rke2-server'
done
