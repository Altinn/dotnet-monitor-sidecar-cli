## dmsctl

CLI to add, remove and connect to dotnet-moniter sidecar in kubernetes

### Synopsis

Small cli tool to ease the process of adding, removing and connecting to dotnet-monitor container to debug a dotnet application running in kubernetes.
Examples:
	# Add sidecars to the pods assosiated with a deployment in kubernetes
	dmsctl add deployment my-deployment

	# Add sidecars to the pods assosiated with a daemonset in kubernetes
	dmsctl add daemonset my-daemonset

	# Setup a port forward from your machine to the dotnet-monitor sidecar container in the pod
	dmsctl port-forward my-pod-13fa7

	# Remove sidecars from the pods assosiated with a deployment in kubernetes
	dmsctl remove deployment my-deployment

	# Remove sidecars from the pods assosiated with a daemonset in kubernetes
	dmsctl remove daemonset my-daemonset 


### Options

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
  -h, --help                help for dmsctl
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl add](dmsctl_add.md)	 - Add a debug sidecar to your pods
* [dmsctl port-forward](dmsctl_port-forward.md)	 - Forward port 52323 from your local machine to port 52323 in a pod
* [dmsctl remove](dmsctl_remove.md)	 - Remove debug sidecar from your pods
* [dmsctl version](dmsctl_version.md)	 - Print the cli version

