#!/usr/bin/env bash

NUMBER_OF_VMS=3

echo "create rke2 config.yaml"

vagrant ssh vm1 -c 'sudo mkdir -p /etc/rancher/rke2'

vagrant ssh vm1 -c '
cat <<EOF > ~/config.yaml
node-name: node-master-01
token:
debug: true
# disable: rke2-ingress-nginx
cni:
  - canal
disable-cloud-controller: true
enable-servicelb: true
kube-apiserver-arg:
  - "default-not-ready-toleration-seconds=30"
  - "default-unreachable-toleration-seconds=30"
EOF
'

vagrant ssh vm1 -c 'sudo mv config.yaml /etc/rancher/rke2'


echo "Enable rke2-server"
vagrant ssh vm1 -c 'sudo systemctl enable rke2-server'
echo "Start rke2-server"
vagrant ssh vm1 -c 'sudo systemctl start rke2-server'

echo "Install kubectl"

kubectl_exists=$(vagrant ssh vm1 -c "command -v kubectl")
if [ -z $kubectl_exists ]; then
  echo "kubectl not installed. Install it"
  vagrant ssh vm1 -c "sudo ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/"
else
  echo "kubectl already installed. skip"
fi
