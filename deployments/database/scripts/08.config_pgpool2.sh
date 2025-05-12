#!/bin/bash

source ./vm_ip.sh

echo "Create log file for pgpool"
sudo touch /var/log/pgpool.log
sudo chown syslog:adm /var/log/pgpool.log  # check permission of other log in /var/log for references
sudo chmod 640 /var/log/pgpool.log

echo "Config pgpool log file"
sudo tee /etc/rsyslog.d/99-pgpool.conf << EOF
local0.*    /var/log/pgpool.log
EOF

echo "Restart rsyslog"
sudo systemctl restart rsyslog

echo "Config log rotation for pgpool"
sudo tee /etc/logrotate.d/pgpool << EOF
/var/log/pgpool.log {      # Target log file to rotate
    weekly                 # Rotate logs weekly
    missingok              # missing ok Don't error if the log file is missing
    rotate 7               # Keep 7 rotated log archives
    compress               # Compress rotated logs (with gzip by default)
    delaycompress          # Wait to compress the previous log until next rotation
    notifempty             # not if empty Don't rotate if the log file is empty
    create 640 syslog adm  # Recreate log file with permissions 640, owner syslog, group adm
    sharedscripts          # Run scripts only once (for multiple log files)
    postrotate             # Command to run after rotation
        systemctl reload rsyslog >/dev/null 2>&1 || true
    endscript
}
EOF

echo "Config pgpool"
sudo tee /etc/pgpool2/pgpool.conf << EOF

#------------------------------------------------------------------------------
# CONNECTION AND AUTHENTICATION
# https://www.pgpool.net/docs/latest/en/html/runtime-config-connection.html
#------------------------------------------------------------------------------
listen_addresses = '*'
port = 9999

# TODO: add config for num_init_children: default 32 and config for authentication

#------------------------------------------------------------------------------
# BACKEND SETTINGS
# https://www.pgpool.net/docs/latest/en/html/runtime-config-backend-settings.html
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
# LOGS
# https://www.pgpool.net/docs/latest/en/html/runtime-config-logging.html
#------------------------------------------------------------------------------
# - Where to log -
log_destination = 'syslog'

# - What to log -
#log_line_prefix = '%m: %a pid %p: '        # Printf-style string to output at beginning of each log line.

log_connections = off
log_disconnections = off
log_pcp_processes = on
log_hostname = off
log_statement = off
log_per_node_statement = off
notice_per_node_statement = off
log_client_messages = off
log_backend_messages = none                 # Log any backend messages. Valid values are none, terse and verbose
log_standby_delay = 'if_over_threshold'     # Log standby delay Valid values are combinations of always, if_over_threshold, none

# - Syslog specific -
syslog_facility = 'LOCAL0'                  # Syslog local facility. Default to LOCAL0
syslog_ident = 'pgpool'                     # Syslog program identification string 

#------------------------------------------------------------------------------
# WATCHDOG CONFIGURATION (pgpool HA)
# https://www.pgpool.net/docs/latest/en/html/runtime-watchdog-config.html
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

wd_lifecheck_method = 'heartbeat'
heartbeat_hostname0 = '${NODE_1_NAME}'
heartbeat_port0 = 9694
heartbeat_device0 = '${DEVICE_INTERFACE}'

heartbeat_hostname1 = '${NODE_2_NAME}'
heartbeat_port1 = 9694
heartbeat_device1 = '${DEVICE_INTERFACE}'

heartbeat_hostname2 = '${NODE_3_NAME}'
heartbeat_port2 = 9694
heartbeat_device2 = '${DEVICE_INTERFACE}'

EOF

echo "Config pgpool node id"
sudo tee /etc/pgpool2/pgpool_node_id << EOF
${PGPOOL_NODE_ID}
EOF

echo "Restart and enable pgpool 2"
sudo systemctl restart pgpool2
sudo systemctl enable pgpool2
