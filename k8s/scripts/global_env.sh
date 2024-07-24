#!/bin/bash

# Define an array of objects
declare -a vms=(
    "192.168.56.101 vm1 master-1"
    "192.168.56.102 vm2 master-2"
    "192.168.56.103 vm3 master-3"
)

RKE2_DOWNLOAD_DIRECTORY=./rke2-artifacts
HELM_CHART_DOWNLOAD_DIRECTORY=./rke2-artifacts/helm-chart
