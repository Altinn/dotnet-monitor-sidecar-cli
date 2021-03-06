## dmsctl remove deployment

Remove a debug sidecar to a Deployment.apps

### Synopsis

After you are done debugging, you can remove the debug sidecar from your pods.
Example:
	# Remove the debug sidecar from a Deployments pods
	dmsctl remove deployment my-deployment

```
dmsctl remove deployment [name] [flags]
```

### Options

```
  -h, --help   help for deployment
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl remove](dmsctl_remove.md)	 - Remove debug sidecar from your pods

