package errors

// IsAlreadyPresent returns true if the error is due to the debug sidecar already being present
func IsAlreadyPresent(err error) bool {
	return err.Error() == "debug sidecar already present"
}

// IsNotPresent returns true if the error is due to the debug sidecar not being present
func IsNotPresent(err error) bool {
	return err.Error() == "debug sidecar not present"
}

// IsNotFound returns true if the error is due to the resource could not be found
func IsNotFound(err error) bool {
	return err.Error() == "resource not found"
}
