//go:build !windows

package state

// sameDir compares two absolute directory paths.
func sameDir(a, b string) bool {
	return a == b
}
