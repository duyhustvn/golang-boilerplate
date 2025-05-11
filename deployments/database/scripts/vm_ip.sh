#!/bin/bash

DEVICE_INTERFACE='enp0s8' # update to your interface

NODE_1_IP=192.168.56.111
NODE_2_IP=192.168.56.112
NODE_3_IP=192.168.56.113
VIP=192.168.56.100

VM_LIST=(
  "$NODE_1_IP postgres01 MASTER"
  "$NODE_2_IP postgres02 BACKUP"
  "$NODE_3_IP postgres03 BACKUP"
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

      if [[ "$ip" == "$vm_ip" ]]; then
        NODE_IP="$ip"
        NODE_NAME="$vm_name"
        KEEPALIVED_ROLE="$keepalived_role"
        if [[ "$KEEPALIVED_ROLE" == "MASTER"  ]]; then
	        PRIORITY=100
	      else
		      PRIORITY=90
        fi
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
