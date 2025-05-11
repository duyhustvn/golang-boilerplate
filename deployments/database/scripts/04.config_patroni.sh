#!/bin/bash

source ./vm_ip.sh

PG_DIR="/u01/data/postgresql"
DATA_DIR="$PG_DIR/16/main"
PG_BIN_DIR="/usr/lib/postgresql/16/bin"
SECRET_DIR="$PG_DIR/16/secret"

echo "Create data folder"
sudo mkdir -p $DATA_DIR
sudo chown -R postgres. $PG_DIR
sudo chmod -R 700 $DATA_DIR

echo "Create secret data folder"
sudo mkdir -p $SECRET_DIR
sudo chmod -R 700 $SECRET_DIR
sudo chown -R postgres. $SECRET_DIR

NAMESPACE="/postgresql-common/" # logical path used within the Distributed Configuration Store (DCS) (e.g., etcd, ZooKeeper) to organize cluster metadata
SCOPE="16-prod"

echo "Config patroni"
sudo tee /etc/patroni/config.yml << EOF
namespace: ${NAMESPACE}
scope: ${SCOPE}
name: ${NODE_NAME}

restapi:
  listen: 0.0.0.0:8008
  connect_address: ${NODE_IP}:8008

etcd3:
  hosts: ${NODE_1_IP}:2379,${NODE_2_IP}:2379,${NODE_3_IP}:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576
    slots:
      repl_slot:
        type: physical

    postgresql:
      use_pg_rewind: true
      use_slots: true
      parameters:
        wal_level: replica
        hot_standby: "on"
        wal_keep_segments: 10
        max_wal_senders: 5
        max_replication_slots: 10
        wal_log_hints: "on"
        logging_collector: 'on'
        max_connections: 512
      pg_hba:
        - local   all             all                                     peer
        - host    all             all             127.0.0.1/32            scram-sha-256
        - host    all             all             ::1/128                 scram-sha-256
        - host    all             all             0.0.0.0/0               scram-sha-256
        - local   replication     all                                     peer
        - host    replication     all             127.0.0.1/32            scram-sha-256
        - host    replication     all             ::1/128                 scram-sha-256
        - host    replication     all             0.0.0.0/0               scram-sha-256

  initdb:
    - encoding: UTF8
    - data-checksums

postgresql:
  listen: 0.0.0.0:5432
  connect_address: ${NODE_IP}:5432
  data_dir: ${DATA_DIR}
  bin_dir: ${PG_BIN_DIR}
  pgpass: ${SECRET_DIR}/pgpass
  authentication:
    replication:
      username: replicator
      password: replPasswd
    superuser:
      username: postgres
      password: postgresPasswd
    rewind:
      username: rewind_user
      password: rewind_password
  parameters:
    unix_socket_directories: "/var/run/postgresql/"
  create_replica_methods:
    - basebackup
  basebackup:
    checkpoint: 'fast'

tags:
  nofailover: false
  noloadbalance: false
  clonefrom: false
  nosync: false
EOF

echo "Restart and enable service"
sudo systemctl restart patroni
sudo systemctl enable patroni

echo "Check status of patroni"
patronictl -c /etc/patroni/config.yml list 16-prod 
