//go:build !windows

package main

// scanOSSpecificPorts is a no-op on non-Windows platforms as the standard
// serial enumerator handles Unix device paths adequately.
func scanOSSpecificPorts(portMap map[string]bool) {
    // No registry to scan on Linux/macOS
}