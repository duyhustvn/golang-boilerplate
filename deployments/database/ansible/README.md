## Vault 
- I already encrypted the information of hosts using this command 
``` sh
ansible-vault encrypt --vault-password-file=.ansible-vault-staging inventory/staging/host_vars/*
```

- To view vault using the command bellow
``` sh
ansible-vault view --vault-password-file=.ansible-vault-staging inventory/staging/host_vars/node-db-01.yml
```

- To decrypt vault 
``` sh
ansible-vault encrypt --vault-password-file=.ansible-vault-staging inventory/staging/host_vars/*
```
You should only decrypt vaul when you want to update, then remember to encrypt it again before commit

## Run ansible 
---
**NOTE**
Please read about (tags)[https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_tags.html] in ansible
---

- Run the playbook
``` sh
ansible-playbook --vault-password-file=.ansible-vault-staging -i inventory/staging/hosts.yml playbooks/install-postgresql-cluster.yml --tags=rc_all
```

