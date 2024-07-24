#!/usr/bin/env bash

source ./global_env.sh

K8S_DASHBOARD_VERSION=7.5.0

if [ ! -f "$HELM_CHART_DOWNLOAD_DIRECTORY/kubernetes-dashboard-$K8S_DASHBOARD_VERSION.tgz" ]; then
    echo "k8s dashboard was not downloaded. Start downloading"
    mkdir -p $HELM_CHART_DOWNLOAD_DIRECTORY
    curl -Lo $HELM_CHART_DOWNLOAD_DIRECTORY/kubernetes-dashboard-$K8S_DASHBOARD_VERSION.tgz https://github.com/kubernetes/dashboard/releases/download/kubernetes-dashboard-$K8S_DASHBOARD_VERSION/kubernetes-dashboard-$K8S_DASHBOARD_VERSION.tgz
else
    echo "k8s dashboard was downloaded. Skip it"
fi


for vm in "${vms[@]}"; do
    ip=$(echo $vm | awk '{print $1}')
    vmname=$(echo $vm | awk '{print $2}')
    hostname=$(echo $vm | awk '{print $3}')

    echo "Copy k8s dashboard to $vmname"
    vagrant rsync $vmname

    echo "Install k8s dashboard to virtual machine $vmname"

    vagrant ssh $vmname -c "sudo helm upgrade --install kubernetes-dashboard /root/rke2-artifacts/helm-chart/kubernetes-dashboard-$K8S_DASHBOARD_VERSION.tgz"
done

# Result of install k8s dashboard
# Release "kubernetes-dashboard" has been upgraded. Happy Helming!
# NAME: kubernetes-dashboard
# LAST DEPLOYED: Wed Jul 24 15:20:49 2024
# NAMESPACE: default
# STATUS: deployed
# REVISION: 3
# TEST SUITE: None
# NOTES:
# *************************************************************************************************
# *** PLEASE BE PATIENT: Kubernetes Dashboard may need a few minutes to get up and become ready ***
# *************************************************************************************************
#
# Congratulations! You have just installed Kubernetes Dashboard in your cluster.
#
# To access Dashboard run:
#   kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443
#
# NOTE: In case port-forward command does not work, make sure that kong service name is correct.
#       Check the services in Kubernetes Dashboard namespace using:
#         kubectl -n default get svc
#
# Dashboard will be available at:
#   https://localhost:8443
