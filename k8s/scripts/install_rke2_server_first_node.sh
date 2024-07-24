#!/usr/bin/env bash

# Include the file that defines the array
source ./global_env.sh

vm=${vms[0]}
echo "vm: $vm"
ip=$(echo $vm | awk '{print $1}')
vmname=$(echo $vm | awk '{print $2}')
hostname=$(echo $vm | awk '{print $3}')

echo "create rke2 config.yaml on $hostname"
vagrant ssh $vmname -c 'sudo mkdir -p /etc/rancher/rke2'
vagrant ssh $vmname -c "
cat <<EOF > ~/config.yaml
node-name: $hostname
node-ip: $ip
token:
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
vagrant ssh $vmname -c 'sudo mv config.yaml /etc/rancher/rke2'


echo "Enable rke2-server"
vagrant ssh $vmname -c 'sudo systemctl enable rke2-server'
echo "Start rke2-server"
vagrant ssh $vmname -c 'sudo systemctl start rke2-server'

echo "Install kubectl"

kubectl_exists=$(vagrant ssh $vmname -c "command -v kubectl")
if [ -z $kubectl_exists ]; then
  echo "kubectl not installed. Install it"
  vagrant ssh $vmname -c "sudo ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/"
else
  echo "kubectl already installed. skip"
fi

echo "Make symbol link for kube config to /root/.kube/config"
vagrant ssh $vmname -c "sudo mkdir -p /root/.kube && sudo ln -s /etc/rancher/rke2/rke2.yaml /root/.kube/config"
