## dmsctl add daemonset

Add the debug sidecar to a daemonsets pods

### Synopsis

To debug your Daemonset, you can add the debug sidecar to your pods.
After you have added the sidecar, you can port forward to one of the pods with dmsctl port-forward [podname].
Example:
	# Add the debug sidecar to a Deployments pods
	dmsctl add daemonset my-daemonset

```
dmsctl add daemonset [name] [flags]
```

### Options

```
  -c, --container string    Supply container name if deployment contains multiple pods
      --debugimage string   image to add as a debug sidecar (default "mcr.microsoft.com/dotnet/monitor:6.0")
  -h, --help                help for daemonset
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.dmsconfig.yaml)
      --kubeconfig string   Override path to the kubeconfig file to use for CLI requests.
  -n, --namespace string    If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.
```

### SEE ALSO

* [dmsctl add](dmsctl_add.md)	 - Add a debug sidecar to your pods

