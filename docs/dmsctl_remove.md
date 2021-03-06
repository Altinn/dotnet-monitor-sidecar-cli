## dmsctl remove

Remove debug sidecar from your pods

### Synopsis

After you are done debugging, you can remove the debug sidecar from your pods.
Example:
	# Remove the debug sidecar from a Deployments pods
	dmsctl remove deployment my-deployment
	# Add the debug sidecar to a DaemonSets pods
	dmsctl remove daemonset my-daemonset

### Options

```
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl](dmsctl.md)	 - CLI to add, remove and connect to dotnet-moniter sidecar in kubernetes
* [dmsctl remove daemonset](dmsctl_remove_daemonset.md)	 - Remove a debug sidecar to a Daemonset.apps
* [dmsctl remove deployment](dmsctl_remove_deployment.md)	 - Remove a debug sidecar to a Deployment.apps

