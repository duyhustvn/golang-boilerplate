#!/usr/bin/env bash

NUMBER_OF_VMS=3

echo "create rke2 config.yaml"

token_path="/var/lib/rancher/rke2/server/token"

token=$(vagrant ssh vm1 -c "sudo cat $token_path")

echo "token: $token"


for ((i = 2; i <= $NUMBER_OF_VMS; i++)); do
  machine=vm$i
  vagrant ssh $machine -c 'sudo mkdir -p /etc/rancher/rke2'

  node_name=node-master-0$i

  vagrant ssh $machine -c "
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
  vagrant ssh $machine -c 'sudo mv config.yaml /etc/rancher/rke2'

  echo "Enable rke2-server"
  vagrant ssh $machine -c 'sudo systemctl enable rke2-server'

  echo "Start rke2-server"
  vagrant ssh $machine -c 'sudo systemctl start rke2-server'

  kubectl_exists=$(vagrant ssh $machine -c "command -v kubectl")
  if [ -z $kubectl_exists ]; then
    echo "kubectl not installed. Install it"
    vagrant ssh $machine -c "sudo ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/"
  else
    echo "kubectl already installed. skip"
  fi
done
