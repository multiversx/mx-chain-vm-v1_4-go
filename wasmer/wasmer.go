// Package wasmer is a Go library to run WebAssembly binaries.
package wasmer

// ForceInstallSighandlers triggers a forced installation of signal handlers in Wasmer 1
func ForceInstallSighandlers() {
	cWasmerForceInstallSighandlers()
}
