package cmd

import (
	"maps"
	"testing"

	"github.com/dittofleet/port-pool/internal/config"
	"github.com/dittofleet/port-pool/internal/state"
)

func TestPerformProvisionRecyclesLeastRecentlyUsedPorts(t *testing.T) {
	cfg := &config.PoolConfig{PortRangeStart: 3100, PortRangeEnd: 3101}
	pcfg := &config.ProjectConfig{PortNames: []string{"web"}}
	s := &state.State{Allocations: []state.Allocation{
		{Dir: "old", Ports: map[string]int{"web": 3100}, Timestamp: 1},
		{Dir: "live", Ports: map[string]int{"web": 3101}, Timestamp: 2},
	}}

	ports, err := PerformProvision(cfg, pcfg, s, "new")
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]int{"web": 3100}
	if !maps.Equal(ports, want) {
		t.Errorf("returned ports = %v, want %v", ports, want)
	}

	wantByDir := map[string]map[string]int{
		"live": {"web": 3101},
		"new":  want,
	}
	if len(s.Allocations) != len(wantByDir) {
		t.Fatalf("allocations = %+v, want dirs live and new", s.Allocations)
	}
	for _, a := range s.Allocations {
		if !maps.Equal(a.Ports, wantByDir[a.Dir]) {
			t.Errorf("%s has ports %v, want %v", a.Dir, a.Ports, wantByDir[a.Dir])
		}
	}
}
