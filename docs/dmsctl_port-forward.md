## dmsctl port-forward

Forward port 52323 from your local machine to port 52323 in a pod

### Synopsis

The debug image does not expose its endpoint out of the pod.
This command will forward port 52323 from your local machine to port 52323 in a pod.
Example:
	# Forward port 52323 from your local machine to port 52323 in the pod my-pod
	dmsctl port-forward my-pod

```
dmsctl port-forward [podname] [flags]
```

### Options

```
  -h, --help   help for port-forward
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl](dmsctl.md)	 - CLI to add, remove and connect to dotnet-moniter sidecar in kubernetes

