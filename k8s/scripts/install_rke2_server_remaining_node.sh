#!/usr/bin/env bash

# Include the file that defines the array
source ./vm_array.sh

token_path="/var/lib/rancher/rke2/server/token"

vm1=${vms[0]}
echo "vm1: $vm1"
vmname1=$(echo $vm1 | awk '{print $2}')
hostname1=$(echo $vm1 | awk '{print $3}')

token=$(vagrant ssh $vmname1 -c "sudo cat $token_path")

echo "token: $token"


# Loop through the array using indices
for (( i=1; i<${#vms[@]}; i++ )); do
    vm=${vms[$i]}
    echo "vm: $vm"
    vmname=$(echo $vm | awk '{print $2}')
    hostname=$(echo $vm | awk '{print $3}')

    vagrant ssh $vmname -c 'sudo mkdir -p /etc/rancher/rke2'
    vagrant ssh $vmname -c "
cat <<EOF > ~/config.yaml
node-name: $hostname
server: https://$hostname1:9345
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
    vagrant ssh $vmname -c 'sudo mv config.yaml /etc/rancher/rke2'

done

#for ((i = 2; i <= $NUMBER_OF_VMS; i++)); do
#  machine=vm$i
#  vagrant ssh $machine -c 'sudo mkdir -p /etc/rancher/rke2'
#
#  node_name=node-master-0$i
#
#  vagrant ssh $machine -c "
#cat <<EOF > ~/config.yaml
#node-name: $node_name
#token: $token
#debug: true
## disable: rke2-ingress-nginx
#cni:
#  - canal
#disable-cloud-controller: true
#enable-servicelb: true
#kube-apiserver-arg:
#  - 'default-not-ready-toleration-seconds=30'
#  - 'default-unreachable-toleration-seconds=30'
#EOF
#"
#  vagrant ssh $machine -c 'sudo mv config.yaml /etc/rancher/rke2'
#
#  echo "Enable rke2-server"
#  vagrant ssh $machine -c 'sudo systemctl enable rke2-server'
#
#  echo "Start rke2-server"
#  vagrant ssh $machine -c 'sudo systemctl start rke2-server'
#
#  kubectl_exists=$(vagrant ssh $machine -c "command -v kubectl")
#  if [ -z $kubectl_exists ]; then
#    echo "kubectl not installed. Install it"
#    vagrant ssh $machine -c "sudo ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/"
#  else
#    echo "kubectl already installed. skip"
#  fi
#done
