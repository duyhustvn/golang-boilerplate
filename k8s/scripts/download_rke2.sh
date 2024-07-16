#!/usr/bin/env bash

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

echo "Copy Rke2 to virtual machine"
for ((i = 1; i <= $NUMBER_OF_VMS; i++)); do
    echo "Copy Rke2 to virtual machine vm$i"
    vagrant rsync "vm$i"
done

echo "Install Rke2 to virtual machine"
for ((i = 1; i <= $NUMBER_OF_VMS; i++)); do
    echo "Install Rke2 to virtual machine vm$i"
    vagrant ssh vm$i -c "sudo INSTALL_RKE2_ARTIFACT_PATH=/root/rke2-artifacts bash /root/rke2-artifacts/install.sh"
done
