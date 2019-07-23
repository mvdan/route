// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package route

import (
	"os/exec"
	"strings"
)

// TODO: Use net.IP? Technically a better type, but clunkier to use, and we have
// to stringify anyway.
// TODO: "result" is probably more than a simple string. Reevaluate once we test
// more platforms and networks.

func to(addr string) (result string, err error) {
	cmd := exec.Command("ip", "route", "get", addr)
	out, err := cmd.CombinedOutput()
	if err, ok := err.(*exec.Error); ok && err.Err == exec.ErrNotFound {
		// "ip route" not available.
		return "", errUnsupportedPlatform{}
	}
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 2 {
		// Network is unreachable.
		return "", nil
	}
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(out))
	for i, field := range fields {
		if field == "via" && i < len(fields)-1 {
			return fields[i+1], nil
		}
	}
	return "", nil
}
