#!/bin/bash

echo "Stop and disable etcd, patroni, postgresql"
sudo systemctl stop {etcd,patroni,postgresql,pgpool2,pgbouncer}
sudo systemctl disable {etcd,patroni,postgresql,pgpool2,pgbouncer}

echo "Delete default value of postgresql"
sudo rm -rf /var/lib/postgresql/16/main
