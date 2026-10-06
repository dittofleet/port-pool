package state

const (
	// minPort and maxPort bound the ports port-pool hands out. Ports below
	// 1024 are well-known ports (and need root on Linux), which a randomly
	// provisioned dev server has no business taking. 65535 is the highest
	// TCP port.
	minPort = 1024
	maxPort = 65535
)

// browserBlockedPorts are the ports browsers refuse to load (Chrome reports
// ERR_UNSAFE_PORT). They are the entries of the "bad port" list in the Fetch
// standard that are at least minPort:
// https://fetch.spec.whatwg.org/#port-blocking
var browserBlockedPorts = map[int]struct{}{
	1719: {}, 1720: {}, 1723: {}, 2049: {}, 3659: {}, 4045: {}, 4190: {},
	5060: {}, 5061: {}, 6000: {}, 6566: {}, 6665: {}, 6666: {}, 6667: {},
	6668: {}, 6669: {}, 6679: {}, 6697: {}, 10080: {},
}

// unusablePort reports whether p is a port port-pool never hands out,
// whatever the pool config says: one outside minPort..maxPort, or one
// browsers refuse to load.
func unusablePort(p int) bool {
	if p < minPort || p > maxPort {
		return true
	}
	_, blocked := browserBlockedPorts[p]
	return blocked
}

// HasUnusablePort reports whether any of the allocation's ports is one
// port-pool no longer hands out. Allocations made by older versions can
// still hold one.
func (a Allocation) HasUnusablePort() bool {
	for _, p := range a.Ports {
		if unusablePort(p) {
			return true
		}
	}
	return false
}
