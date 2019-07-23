// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package route

func IsUnsupported(err error) bool {
	return err == errUnsupportedPlatform{}
}

type errUnsupportedPlatform struct{}

func (errUnsupportedPlatform) Error() string {
	return "unsupported platform"
}

// Offline returns whether the machine is known to be offline. The function
// returns true if there's no valid route to 8.8.8.8. Unsupported platforms
// always return false.
func Offline() bool {
	via, err := to("8.8.8.8")
	if IsUnsupported(err) {
		return false
	}
	if err != nil {
		return true
	}
	return via == ""
}
