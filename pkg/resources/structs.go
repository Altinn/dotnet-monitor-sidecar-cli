package resources

// DDConfig represents the configuration applied for the debug sidecar
type DDConfig struct {
	// The name of the container to debug
	ContainerToDebug string `json:"containerToDebug"`
	// The name of the sidecar.
	DebugContainerName string `json:"debugContainerName"`
	// TmpdirAdded reflects if tmpdir was added by dd
	TmpdirAdded bool `json:"tmpdirAdded"`
	// SecretMount reflect name of created and mounted
	SecretName string `json:"secretMount"`
}
