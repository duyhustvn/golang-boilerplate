#!/bin/bash

DEVICE_INTERFACE='enp0s8' # update to your interface

NODE_1_IP=192.168.56.111
NODE_2_IP=192.168.56.112
NODE_3_IP=192.168.56.113

NODE_1_NAME=node-db-01
NODE_2_NAME=node-db-02
NODE_3_NAME=node-db-03

ETCD_1_NAME=postgres01
ETCD_2_NAME=postgres02
ETCD_3_NAME=postgres03

VIP=192.168.56.100 

VM_LIST=(
  "$NODE_1_IP $NODE_1_NAME MASTER 0 $ETCD_1_NAME"
  "$NODE_2_IP $NODE_2_NAME BACKUP 1 $ETCD_2_NAME"
  "$NODE_3_IP $NODE_3_NAME BACKUP 2 $ETCD_3_NAME"
)

get_node_ip_from_interface() {
  local interface="$1"

  # Get all IPs from the interface (IPv4 only)
  local ips
  ips=$(ip -4 addr show "$interface" | awk '/inet / {print $2}' | cut -d/ -f1)

  # Loop through each IP on the interface
  for ip in $ips; do
    # Loop through each VM entry
    for vm in "${VM_LIST[@]}"; do
      vm_ip=$(echo "$vm" | awk '{print $1}')
      vm_name=$(echo "$vm" | awk '{print $2}')
      keepalived_role=$(echo "$vm" | awk '{print $3}')
      pgpool_node_id=$(echo "$vm" | awk '{print $4}')
      etcd_name=$(echo "$vm" | awk '{print $5}')

      if [[ "$ip" == "$vm_ip" ]]; then
        NODE_IP="$ip"
        NODE_NAME="$vm_name"
        KEEPALIVED_ROLE="$keepalived_role"
        if [[ "$KEEPALIVED_ROLE" == "MASTER"  ]]; then
	        PRIORITY=100
	      else
		      PRIORITY=90
        fi
        PGPOOL_NODE_ID=$pgpool_node_id
        ETCD_NAME=$etcd_name
        return 0
      fi
    done
  done

  echo "ERROR: No matching NODE_IP found on interface $interface." >&2
  return 1
}

# Call the function
get_node_ip_from_interface "$DEVICE_INTERFACE" || exit 1

# Output result
echo "Matched NODE_IP: $NODE_IP"
echo "Matched NODE_NAME: $NODE_NAME"
echo "Matched KEEPALIVED_ROLE: $KEEPALIVED_ROLE"
echo "Matched PRIORITY: $PRIORITY"
echo "Matched PGPOOL_NODE_ID: $PGPOOL_NODE_ID"
echo "Matched ETCD_NAME: $ETCD_NAME"
