package vmmd

import (
	"os"
)

// FIPSEnabled returns true if running in FIPS mode.
// currently it just checks the system wide /etc/system-fips file present or not.
// We can improve it later.
func FIPSEnabled() bool {
	fips_file := "/etc/system-fips"
	if _, err := os.Stat(fips_file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
