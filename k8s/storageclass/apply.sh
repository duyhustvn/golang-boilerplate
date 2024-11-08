#!/usr/bin/env bash

kubectl apply -f storageclass.yaml
# Set local-path-storage to default storageclass
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
