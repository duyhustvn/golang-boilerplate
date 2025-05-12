#!/bin/bash

source ./vm_ip.sh

echo "Config pgpool"
sudo tee /etc/pgpool2/pgpool.conf << EOF
listen_addresses = '*'
port = 9999

#------------------------------------------------------------------------------
# BACKEND CONFIGURATION (PostgreSQL Nodes)
#------------------------------------------------------------------------------
backend_hostname0 = ${NODE_1_NAME}
backend_port0 = 5432
backend_data_directory0 = '/u01/data/postgresql/16/main'
backend_flag0 = 'ALLOW_TO_FAILOVER'

backend_hostname1 = ${NODE_2_NAME}
backend_port1 = 5432
backend_data_directory1 = '/u01/data/postgresql/16/main'
backend_flag1 = 'ALLOW_TO_FAILOVER'

backend_hostname2 = ${NODE_3_NAME}
backend_port2 = 5432
backend_data_directory2 = '/u01/data/postgresql/16/main'
backend_flag2 = 'ALLOW_TO_FAILOVER'

#------------------------------------------------------------------------------
# WATCHDOG CONFIGURATION (pgpool HA)
#------------------------------------------------------------------------------
use_watchdog = on
hostname0 = '${NODE_1_NAME}'  
wd_port0 = 9000
pgpool_port0 = 9999

hostname1 = '${NODE_2_NAME}'  
wd_port1 = 9000
pgpool_port1 = 9999

hostname2 = '${NODE_3_NAME}'  
wd_port2 = 9000
pgpool_port2 = 9999

delegate_ip = '${VIP}'  # Virtual IP

EOF

echo "Config pgpool node id"
sudo tee /etc/pgpool2/pgpool_node_id << EOF
${PGPOOL_NODE_ID}
EOF

echo "Restart and enable pgpool 2"
sudo systemctl restart pgpool2
sudo systemctl enable pgpool2
