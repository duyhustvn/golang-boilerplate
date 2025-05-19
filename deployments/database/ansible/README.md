## Vault 
- I already encrypted the information of hosts using this command 
``` sh
ansible-vault encrypt --vault-password-file=.ansible-vault-staging inventory/staging/host_vars/*
```
- To view vault content for example
``` sh
ansible-vault view --vault-password-file=.ansible-vault-staging inventory/staging/host_vars/node-db-01.yml
```
- To view vault content 

## Run ansible 
ansible-playbook --vault-password-file=.ansible-vault-staging -i inventory/staging/hosts.yml playbooks/install-postgresql-cluster.yml
