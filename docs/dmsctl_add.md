## dmsctl add

Add a debug sidecar to your pods

### Synopsis

To debug your pods, you can add the debug sidecar to your pods.
Example:
	# Add the debug sidecar to a Deployments pods
	dmsctl add deployment my-deployment
	# Add the debug sidecar to a DaemonSets pods
	dmsctl add daemonset my-daemonset

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl](dmsctl.md)	 - CLI to add, remove and connect to dotnet-moniter sidecar in kubernetes
* [dmsctl add daemonset](dmsctl_add_daemonset.md)	 - Add the debug sidecar to a daemonsets pods
* [dmsctl add deployment](dmsctl_add_deployment.md)	 - Add the debug sidecar to a deployments pods

