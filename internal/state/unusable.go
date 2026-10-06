package state

import "github.com/dittofleet/port-pool/internal/config"

// browserBlockedPorts are the ports browsers refuse to load (Chrome reports
// ERR_UNSAFE_PORT), taken from the "bad port" list in the Fetch standard:
// https://fetch.spec.whatwg.org/#port-blocking
var browserBlockedPorts = map[int]struct{}{
	0: {}, 1: {}, 7: {}, 9: {}, 11: {}, 13: {}, 15: {}, 17: {}, 19: {}, 20: {},
	21: {}, 22: {}, 23: {}, 25: {}, 37: {}, 42: {}, 43: {}, 53: {}, 69: {},
	77: {}, 79: {}, 87: {}, 95: {}, 101: {}, 102: {}, 103: {}, 104: {}, 109: {},
	110: {}, 111: {}, 113: {}, 115: {}, 117: {}, 119: {}, 123: {}, 135: {},
	137: {}, 139: {}, 143: {}, 161: {}, 179: {}, 389: {}, 427: {}, 465: {},
	512: {}, 513: {}, 514: {}, 515: {}, 526: {}, 530: {}, 531: {}, 532: {},
	540: {}, 548: {}, 554: {}, 556: {}, 563: {}, 587: {}, 601: {}, 636: {},
	989: {}, 990: {}, 993: {}, 995: {}, 1719: {}, 1720: {}, 1723: {}, 2049: {},
	3659: {}, 4045: {}, 4190: {}, 5060: {}, 5061: {}, 6000: {}, 6566: {},
	6665: {}, 6666: {}, 6667: {}, 6668: {}, 6669: {}, 6679: {}, 6697: {},
	10080: {},
}

// unusablePort reports whether p is a port port-pool never hands out,
// whatever the pool config says: one outside config.MinPort..MaxPort, or one
// browsers refuse to load.
func unusablePort(p int) bool {
	if p < config.MinPort || p > config.MaxPort {
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
