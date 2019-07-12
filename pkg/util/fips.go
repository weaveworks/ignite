package util

import (
	"os"
)

// FIPSEnabled returns true if running in FIPS mode.
// currently it just checks the system wide /etc/system-fips file present or not.
// TODO - Find a better generic solution for this.
func FIPSEnabled() bool {
	fips_file := "/etc/system-fips"
	if _, err := os.Stat(fips_file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
