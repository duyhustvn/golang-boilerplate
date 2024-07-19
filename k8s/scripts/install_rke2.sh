#!/usr/bin/env bash

# Include the file that defines the array
source ./vm_array.sh

RKE2_VERSION=v1.29.4+rke2r1
RKE2_DOWNLOAD_DIRECTORY=./rke2-artifacts
NUMBER_OF_VMS=3

echo "Check if Rke2 was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/rke2-images.linux-amd64.tar.zst" ]; then
    echo "Rke2 was not downloaded. Start downloading"
    wget -O $RKE2_DOWNLOAD_DIRECTORY/rke2-images.linux-amd64.tar.zst https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/rke2-images.linux-amd64.tar.zst
else
    echo "Rke2 already downloaded"
fi

echo "Check if Rke2 binary was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/rke2.linux-amd64.tar.gz" ]; then
    echo "Rke2 binary was not downloaded. Start downloading"
    wget -O $RKE2_DOWNLOAD_DIRECTORY/rke2.linux-amd64.tar.gz https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/rke2.linux-amd64.tar.gz
else
    echo "Rke2 binary already downloaded"
fi

echo "Check if Rke2 checksum was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/sha256sum-amd64.txt" ]; then
    echo "Rke2 checksum was not downloaded. Start downloading"
    wget -O $RKE2_DOWNLOAD_DIRECTORY/sha256sum-amd64.txt https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/sha256sum-amd64.txt
else
    echo "Rke2 checksum already downloaded"
fi

echo "Check if Rke2 install script was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/install.sh" ]; then
    echo "Rke2 install script was not downloaded. Start downloading"
    wget -O $RKE2_DOWNLOAD_DIRECTORY/install.sh https://get.rke2.io
else
    echo "Rke2 install script downloaded"
fi

# Loop through the array and add entries to the /etc/hosts file
for vm in "${vms[@]}"; do
    ip=$(echo $vm | awk '{print $1}')
    vmname=$(echo $vm | awk '{print $2}')
    echo "Copy Rke2 to virtual machine $vmname"
    vagrant rsync $vmname

    echo "Install Rke2 to virtual machine $vmname"
    rke2_installed=$(vagrant ssh $vmname -c "command -v rke2")
    if [ -z $rke2_installed ]; then
        echo "Rke2 not installed install it"
        vagrant ssh $vmname -c "sudo INSTALL_RKE2_ARTIFACT_PATH=/root/rke2-artifacts bash /root/rke2-artifacts/install.sh"
    else
        echo "Rke2 already installed skip"
    fi
done
