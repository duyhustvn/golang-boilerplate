#!/bin/bash

echo "Install Postgresql 16"
# Import the repository signing key:
sudo apt install curl ca-certificates
sudo install -d /usr/share/postgresql-common/pgdg
sudo curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc

# Create the repository configuration file:
. /etc/os-release
sudo sh -c "echo 'deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc]  https://apt.postgresql.org/pub/repos/apt $VERSION_CODENAME-pgdg main' > /etc/apt/sources.list.d/pgdg.list"
sudo apt update
sudo apt install -y postgresql-16

echo "Install patroni and etcd"
sudo apt install -y python3-pip python3-dev binutils 
sudo apt install -y patroni etcd etcd-server etcd-client pgbackrest

echo "Install HAProxy and Keepalived"
sudo apt install -y haproxy keepalived

echo "Install Pgbouncer"
sudo apt install -y pgbouncer

echo "Install Pgpool 2"
sudo apt install -y pgpool2
