package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePoolConfig(t *testing.T, contents string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	path := filepath.Join(dir, "port-pool", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadPoolDefaultsPortRange(t *testing.T) {
	cases := []struct {
		config     string
		start, end int
	}{
		{`{"schemaVersion": 1}`, DefaultPortRangeStart, DefaultPortRangeEnd},
		{`{"schemaVersion": 1, "portRangeStart": 4000}`, 4000, DefaultPortRangeEnd},
		{`{"schemaVersion": 1, "portRangeEnd": 5000}`, DefaultPortRangeStart, 5000},
		{`{"schemaVersion": 1, "portRangeStart": 20000, "portRangeEnd": 30000}`, 20000, 30000},
	}
	for _, c := range cases {
		writePoolConfig(t, c.config)
		cfg, err := LoadPool()
		if err != nil {
			t.Fatalf("%s: %v", c.config, err)
		}
		if cfg.PortRangeStart != c.start || cfg.PortRangeEnd != c.end {
			t.Errorf("%s: got %d-%d, want %d-%d", c.config, cfg.PortRangeStart, cfg.PortRangeEnd, c.start, c.end)
		}
	}
}

func TestLoadPoolRejectsInvalidPortRange(t *testing.T) {
	cases := map[string]string{
		`{"schemaVersion": 1, "portRangeStart": 0}`:     "portRangeStart: must be a positive integer",
		`{"schemaVersion": 1, "portRangeStart": 10000}`: "portRangeEnd (9999) must be >= portRangeStart (10000)",
	}
	for config, want := range cases {
		writePoolConfig(t, config)
		_, err := LoadPool()
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want an error containing %q", config, err, want)
		}
	}
}
