package state

import "testing"

func TestFindByDirIgnoresCase(t *testing.T) {
	s := &State{Allocations: []Allocation{{Dir: `F:\Code\app`}}}
	if a, i := s.FindByDir(`f:\code\APP`); a == nil || i != 0 {
		t.Fatalf("expected the allocation at index 0, got %v at %d", a, i)
	}
}
