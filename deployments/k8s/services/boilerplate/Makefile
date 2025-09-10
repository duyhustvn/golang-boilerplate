apply: 
	sops -d encrypted-secret-dev.yaml > unencrypted-dev-secrets.yaml && helm upgrade --install boilerplate-be-dev --values unencrypted-dev-secrets.yaml ./ && rm unencrypted-dev-secrets.yaml

undeploy: 
	helm uninstall boilerplate-be-dev
