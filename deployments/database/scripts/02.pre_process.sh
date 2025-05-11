#!/bin/bash

echo "Stop and disable etcd, patroni, postgresql"
sudo systemctl stop {etcd,patroni,postgresql}
sudo systemctl disable {etcd,patroni,postgresql}

echo "Delete default value of postgresql"
sudo rm -rf /var/lib/postgresql/16/main
