package state

import "strings"

// sameDir compares two absolute directory paths. Windows paths ignore case,
// and tools disagree on even the drive letter's.
func sameDir(a, b string) bool {
	return strings.EqualFold(a, b)
}
