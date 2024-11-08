#!/usr/bin/env bash


kubectl apply -f rke2-cannal-config.yaml
# Restart all rke2-canal
kubectl delete pod -l k8s-app=canal -n kube-system
