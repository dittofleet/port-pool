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

func TestFindNextAvailablePortsStaysWithin1024To65535(t *testing.T) {
	cases := []struct {
		start, end, want int
	}{
		{1, 1024, 1024},
		{65535, 70000, 65535},
	}
	for _, c := range cases {
		cfg := &config.PoolConfig{PortRangeStart: c.start, PortRangeEnd: c.end}
		ports := FindNextAvailablePorts(&State{}, cfg, 1)
		if len(ports) != 1 || ports[0] != c.want {
			t.Errorf("range %d-%d: got %v, want [%d]", c.start, c.end, ports, c.want)
		}
	}

	cfg := &config.PoolConfig{PortRangeStart: 1, PortRangeEnd: 1023}
	if ports := FindNextAvailablePorts(&State{}, cfg, 1); ports != nil {
		t.Errorf("range 1-1023: got %v, want nil", ports)
	}
}

func TestHasUnusablePort(t *testing.T) {
	for p, want := range map[int]bool{80: true, 1023: true, 1024: false, 6566: true, 8080: false, 65535: false, 65536: true} {
		a := Allocation{Ports: map[string]int{"web": p}}
		if got := a.HasUnusablePort(); got != want {
			t.Errorf("port %d: got %v, want %v", p, got, want)
		}
	}
}
