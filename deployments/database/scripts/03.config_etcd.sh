#!/bin/bash

source ./vm_ip.sh

echo "Create data forder for etcd"
sudo mkdir -p /u01/data/etcd 
sudo chown -R etcd. /u01/data/etcd
sudo chmod -R 700 /u01/data/etcd

echo "Config etcd at /etc/default/etcd"
sudo tee /etc/default/etcd << EOF
ETCD_NAME=${NODE_NAME}
ETCD_DATA_DIR=/u01/data/etcd/postgresql
ETCD_LISTEN_PEER_URLS="http://${NODE_IP}:2380"
ETCD_LISTEN_CLIENT_URLS="http://${NODE_IP}:2379,http://localhost:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://${NODE_IP}:2380"
ETCD_INITIAL_CLUSTER="postgres01=http://${NODE_1_IP}:2380,postgres02=http://${NODE_2_IP}:2380,postgres03=http://${NODE_3_IP}:2380"
ETCD_INITIAL_CLUSTER_STATE="new"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_ADVERTISE_CLIENT_URLS="http://${NODE_IP}:2379"
EOF

echo "Start and enable etcd"
sudo systemctl start etcd
sudo systemctl enable etcd
