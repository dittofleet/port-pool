package state

import (
	"testing"

	"github.com/dittofleet/port-pool/internal/config"
)

func TestFindNextAvailablePortsSkipsBrowserBlockedPorts(t *testing.T) {
	// 6566 sits in the middle of the range, so a block of 2 only fits on
	// either side of it and a block of 3 doesn't fit at all.
	cfg := &config.PoolConfig{PortRangeStart: 6564, PortRangeEnd: 6568}
	bases := map[int]bool{}
	for i := 0; i < 100; i++ {
		ports := FindNextAvailablePorts(&State{}, cfg, 2)
		if ports == nil {
			t.Fatal("expected a block, got nil")
		}
		bases[ports[0]] = true
	}
	if len(bases) != 2 || !bases[6564] || !bases[6567] {
		t.Fatalf("expected blocks to start at 6564 and 6567, got %v", bases)
	}

	if ports := FindNextAvailablePorts(&State{}, cfg, 3); ports != nil {
		t.Fatalf("expected no block of 3 around 6566, got %v", ports)
	}
}
