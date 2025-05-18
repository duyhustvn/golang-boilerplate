#!/bin/bash

source ./vm_ip.sh

sudo -u postgres psql -c "CREATE ROLE pgpool WITH LOGIN PASSWORD 'changeme';"
sudo -u postgres psql -c "CREATE ROLE repl   WITH LOGIN PASSWORD 'changeme';"
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'changeme';"
sudo -u postgres psql -c "GRANT pg_monitor TO pgpool"

sudo cp postgresql/postgresql.conf /etc/postgresql/17/main/postgresql.conf

echo "Config pg_hba for postgresql"
sudo cp postgresql/pg_hba.conf /etc/postgresql/17/main/pg_hba.conf
sudo systemctl restart postgresql.service

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

echo "Create .pgpass"
sudo tee /var/lib/postgresql/.pgpass << EOF
# host:port:database:username:password
${NODE_1_NAME}:5432:replication:repl:changeme
${NODE_2_NAME}:5432:replication:repl:changeme
${NODE_3_NAME}:5432:replication:repl:changeme
${NODE_1_NAME}:5432:postgres:repl:changeme
${NODE_2_NAME}:5432:postgres:repl:changeme
${NODE_3_NAME}:5432:postgres:repl:changeme
EOF

sudo -i -u postgres bash -c "chmod 600 ~/.pgpass"

sudo tee /etc/pgpool2/pcp.conf << EOF
# PCP Client Authentication Configuration File
# ============================================
#
# This file contains user ID and his password for pgpool
# communication manager authentication.
#
# Note that users defined here do not need to be PostgreSQL
# users. These users are authorized ONLY for pgpool
# communication manager.
#
# File Format
# ===========
#
# List one UserID and password on a single line. They must
# be concatenated together using ':' (colon) between them.
# No spaces or tabs are allowed anywhere in the line.
#
# Example:
# postgres:e8a48653851e28c69d0506508fb27fc5
#
# Be aware that there will be no spaces or tabs at the
# beginning of the line! although the above example looks
# like so.
#
# Lines beginning with '#' (pound) are comments and will
# be ignored. Again, no spaces or tabs allowed before '#'.

# USERID:MD5PASSWD
pgpool:$(pg_md5 pgpool_password)
EOF

sudo -i -u postgres bash -c "echo 'localhost:9898:pgpool:pgpool_password' > ~/.pcppass"
sudo -i -u postgres bash -c "chmod 600 ~/.pcppass"

echo "Config pgpool node id"
sudo tee /etc/pgpool2/pgpool_node_id << EOF
${PGPOOL_NODE_ID}
EOF

sudo -i -u postgres bash -c "echo 'secret_string' > ~/.pgpoolkey"
sudo -i -u postgres bash -c "chmod 600 ~/.pgpoolkey"

echo "Enter password for pgpool user in postgresql"
sudo pg_enc -m -k /var/lib/postgresql/.pgpoolkey -u pgpool -p
echo "Enter password for postgres user in postgresql"
sudo pg_enc -m -k /var/lib/postgresql/.pgpoolkey -u postgres -p

sudo cp pgpool2/pool_hba.conf /etc/pgpool2/pool_hba.conf

sudo tee /etc/pgpool2/pgpool.conf << EOF
#------------------------------------------------------------------------------
# CONNECTION AND AUTHENTICATION
# https://www.pgpool.net/docs/latest/en/html/runtime-config-connection.html
#------------------------------------------------------------------------------
backend_clustering_mode = 'streaming_replication'
listen_addresses = '*'
pcp_listen_addresses = '*'
port = 9999

enable_pool_hba = true

#------------------------------------------------------------------------------
# STREAMING REPLICATION MODE
#------------------------------------------------------------------------------
sr_check_user = 'pgpool' # sr: stream replication
sr_check_password = ''

#------------------------------------------------------------------------------
# Health Check
# https://www.pgpool.net/docs/latest/en/html/runtime-config-health-check.html
#------------------------------------------------------------------------------
health_check_period = 5
health_check_timeout = 30
health_check_user = 'pgpool'
health_check_password = ''
health_check_max_retries = 3

#------------------------------------------------------------------------------
# BACKEND SETTINGS
# https://www.pgpool.net/docs/latest/en/html/runtime-config-backend-settings.html
#------------------------------------------------------------------------------
backend_hostname0 = '${NODE_1_NAME}'
backend_port0 = 5432
backend_weight0 = 1
backend_data_directory0 = '/var/lib/postgresql/17/main'
backend_flag0 = 'ALLOW_TO_FAILOVER'

backend_hostname1 = '${NODE_2_NAME}'
backend_port1 = 5432
backend_weight1 = 1
backend_data_directory1 = '/var/lib/postgresql/17/main/'
backend_flag1 = 'ALLOW_TO_FAILOVER'

backend_hostname2 = '${NODE_3_NAME}'
backend_port2 = 5432
backend_weight2 = 1
backend_data_directory2 = '/var/lib/postgresql/17/main/'
backend_flag2 = 'ALLOW_TO_FAILOVER'

#------------------------------------------------------------------------------
# WATCHDOG CONFIGURATION
# https://www.pgpool.net/docs/latest/en/html/runtime-watchdog-config.html
#------------------------------------------------------------------------------
use_watchdog = on
delegate_IP = '${VIP}'  # Virtual IP
if_cmd_path = '/sbin' # path to the directory where if_up/down_cmd exists
if_up_cmd = '/usr/bin/sudo /sbin/ip addr add $_IP_$/24 dev ${DEVICE_INTERFACE} label ${DEVICE_INTERFACE}:0' # startup delegate IP command
if_down_cmd = '/usr/bin/sudo /sbin/ip addr del $_IP_$/24 dev ${DEVICE_INTERFACE}' # shutdown delegate IP command
arping_cmd = '/usr/bin/sudo /usr/sbin/arping -U $_IP_$ -w 1 -I ${DEVICE_INTERFACE}'

hostname0 = '${NODE_1_NAME}'
wd_port0 = 9000
pgpool_port0 = 9999

hostname1 = '${NODE_2_NAME}'
wd_port1 = 9000
pgpool_port1 = 9999

hostname2 = '${NODE_3_NAME}'
wd_port2 = 9000
pgpool_port2 = 9999

wd_lifecheck_method = 'heartbeat'
wd_interval = 10

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

wd_heartbeat_keepalive = 2
wd_heartbeat_deadtime = 30

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

EOF

echo "Restart and enable pgpool 2"
sudo systemctl restart pgpool2
sudo systemctl enable pgpool2
