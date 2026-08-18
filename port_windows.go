//go:build windows

package main

import (
    "golang.org/x/sys/windows/registry"
)

// scanOSSpecificPorts checks the Windows registry for available COM ports
func scanOSSpecificPorts(portMap map[string]bool) {
    key, err := registry.OpenKey(
        registry.LOCAL_MACHINE,
        `HARDWARE\DEVICEMAP\SERIALCOMM`,
        registry.QUERY_VALUE,
    )
    if err != nil {
        return
    }
    defer key.Close()

    names, err := key.ReadValueNames(0)
    if err != nil {
        return
    }

    for _, name := range names {
        value, _, err := key.GetStringValue(name)
        if err == nil && value != "" {
            portMap[value] = true
        }
    }
}