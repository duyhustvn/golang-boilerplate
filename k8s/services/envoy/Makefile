service=eg

apply: 
	helm install $(service) oci://docker.io/envoyproxy/gateway-helm --version v1.5.0 -f helm/values.yml -n envoy-gateway-system --create-namespace

undeploy: 
	helm uninstall $(service)
