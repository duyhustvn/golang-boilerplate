#!/bin/bash

source ./vm_ip.sh

echo "Config keepalived"

sudo tee /etc/keepalived/keepalived.conf << EOF
global_defs {
}

vrrp_script chk_haproxy { # Requires keepalived-1.1.13
    script "killall -0 haproxy" # widely used idiom
    interval 2 # check every 2 seconds
    weight 20 # add 2 points of prio if OK
}

vrrp_instance VI_PSQL {
    interface ${DEVICE_INTERFACE}
    state ${KEEPALIVED_ROLE} 
    priority ${PRIORITY} # 100 on master, 90 on backup
    virtual_router_id 51
    authentication {
        auth_type PASS
        auth_pass changeme
    }
    virtual_ipaddress {
        ${VIP}
    }
    track_script {
        chk_haproxy
    }
}
EOF

echo "Start and Enable keepalived"
sudo systemctl restart keepalived
sudo systemctl enable keepalived
