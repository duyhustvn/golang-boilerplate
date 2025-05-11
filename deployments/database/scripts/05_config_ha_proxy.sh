#!/bin/bash

source ./vm_ip.sh

echo "Config HAProxy"
sudo tee /etc/haproxy/haproxy.cfg << EOF
global
    maxconn 500

defaults
    log global
    mode tcp              # TCP mode (PostgreSQL uses TCP connections)
    retries 2           
    timeout client 30m    # Close client connections after 30m of inactivity
    timeout connect 5s    # Max time to connect to a backend server
    timeout server 30m    # Abort queries if PostgreSQL takes >30m to respond
    timeout check 5s

listen stats
    mode http
    bind *:8404
    stats enable
    stats uri /stats

frontend pg_write
		bind *:5000 # Apps send writes to port 5000
		default_backend pg_primary

frontend pg_read
		bind *:5001 # Apps send reads to port 5001
		default_backend pg_replicas

backend pg_primary
    option httpchk GET /primary  # Health check using Patroni API
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server postgres1 ${NODE_1_IP}:6432 maxconn 100 check port 8008
    server postgres2 ${NODE_2_IP}:6432 maxconn 100 check port 8008
    server postgres3 ${NODE_3_IP}:6432 maxconn 100 check port 8008

backend pg_replicas
    balance roundrobin 
    option httpchk GET /replica # Health check using Patroni API
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server postgres1 ${NODE_1_IP}:6432 maxconn 100 check port 8008
    server postgres2 ${NODE_2_IP}:6432 maxconn 100 check port 8008
    server postgres3 ${NODE_3_IP}:6432 maxconn 100 check port 8008
EOF

echo "Start and enable haproxy"
sudo systemctl restart haproxy
sudo systemctl enable haproxy
