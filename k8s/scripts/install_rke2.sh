#!/usr/bin/env bash

# Include the file that defines the array
source ./global_env.sh

RKE2_VERSION=v1.30.1+rke2r1
HELM_VERSION=v3.14.3

echo "Check if Rke2 was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/rke2-images.linux-amd64.tar.zst" ]; then
    echo "Rke2 was not downloaded. Start downloading https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/rke2-images.linux-amd64.tar.zst"
    curl -o $RKE2_DOWNLOAD_DIRECTORY/rke2-images.linux-amd64.tar.zst -L https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/rke2-images.linux-amd64.tar.zst
else
    echo "Rke2 already downloaded. Skip"
fi

echo "Check if Rke2 binary was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/rke2.linux-amd64.tar.gz" ]; then
    echo "Rke2 binary was not downloaded. Start downloading"
    curl -o $RKE2_DOWNLOAD_DIRECTORY/rke2.linux-amd64.tar.gz -L https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/rke2.linux-amd64.tar.gz
else
    echo "Rke2 binary already downloaded. Skip"
fi

echo "Check if Rke2 checksum was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/sha256sum-amd64.txt" ]; then
    echo "Rke2 checksum was not downloaded. Start downloading"
    curl -o $RKE2_DOWNLOAD_DIRECTORY/sha256sum-amd64.txt -L https://github.com/rancher/rke2/releases/download/$RKE2_VERSION/sha256sum-amd64.txt
else
    echo "Rke2 checksum already downloaded. Skip"
fi

echo "Check if Rke2 install script was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/install.sh" ]; then
    echo "Rke2 install script was not downloaded. Start downloading"
    curl -o $RKE2_DOWNLOAD_DIRECTORY/install.sh -L https://get.rke2.io
else
    echo "Rke2 install script downloaded. Skip"
fi

echo "Check if helm was downloaed"
if [ ! -f "$RKE2_DOWNLOAD_DIRECTORY/helm" ]; then
    echo "Helm was not downloaded. Start downloading"
    curl -o - https://get.helm.sh/helm-$HELM_VERSION-linux-amd64.tar.gz | tar -xz -C $RKE2_DOWNLOAD_DIRECTORY
    mv $RKE2_DOWNLOAD_DIRECTORY/linux-amd64/helm $RKE2_DOWNLOAD_DIRECTORY/helm
    rm -r $RKE2_DOWNLOAD_DIRECTORY/linux-amd64
else
    echo "Helm was downloaded. Skip"
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


    helm_installed=$(vagrant ssh $vmname -c "command -v helm")
    if [ -z $helm_installed ]; then
        echo "Helm not installed. Install it"
        vagrant ssh $vmname -c "sudo cp /root/rke2-artifacts/helm /usr/local/bin/helm"
    else
        echo "Helm already installed. Skip"
    fi
done
