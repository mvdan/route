// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

// +build !linux

package route

func to(addr string) (result string, err error) {
	return "", errUnsupportedPlatform{}
}
