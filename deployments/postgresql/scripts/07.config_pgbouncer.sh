#!/bin/bash

echo "Config pgbouncer"

sudo tee /etc/pgbouncer/pgbouncer.ini << EOF

[databases]
* = host=localhost port=5432 auth_user=postgres

[pgbouncer]
listen_addr = *
listen_port = 6432

auth_type = scram-sha-256
auth_file = /etc/pgbouncer/userlist.txt

pool_mode = transaction
max_client_conn = 10000
default_pool_size = 50
EOF

echo "Restart and enable pgbouncer"
sudo systemctl restart pgbouncer 
sudo systemctl enable pgbouncer
