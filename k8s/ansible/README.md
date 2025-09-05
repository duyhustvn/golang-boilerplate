## Install k8s cluster
```
ansible-playbook --vault-password-file=.ansible-vault-secret-devlocal -i inventory/devlocal/ubuntu/hosts.yml playbooks/install-k8s-cluster.yml
```

