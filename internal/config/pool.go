package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dittofleet/go-cli-kit/xdg"
	"github.com/dittofleet/port-pool/internal/app"
)

const PoolSchemaVersion = 1

// The pool range used when the config leaves portRangeStart or
// portRangeEnd out.
const (
	DefaultPortRangeStart = 3000
	DefaultPortRangeEnd   = 9999
)

// MinPort and MaxPort bound the ports port-pool hands out, whatever the
// configured range. Ports below 1024 are well-known ports (and need root on
// Linux), which a randomly provisioned dev server has no business taking.
// 65535 is the highest TCP port.
const (
	MinPort = 1024
	MaxPort = 65535
)

type PoolConfig struct {
	SchemaVersion  int   `json:"schemaVersion"`
	PortRangeStart int   `json:"portRangeStart"`
	PortRangeEnd   int   `json:"portRangeEnd"`
	ExcludedPorts  []int `json:"excludedPorts"`
}

func PoolConfigPath() string {
	return filepath.Join(xdg.ConfigDir(app.Name), "config.json")
}

// StarterPoolConfig returns the recommended starter content for a fresh
// install. main prints this when LoadPool returns a MissingPoolConfigError.
const StarterPoolConfig = `{
  "schemaVersion": 1,
  "excludedPorts": [
    3000, 3001, 3306,
    4000, 4200,
    5000, 5173, 5432, 5500,
    6379,
    8000, 8080, 8081, 8443, 8888,
    9000, 9090, 9200
  ]
}`

// MissingPoolConfigError is returned by LoadPool when the config file
// doesn't exist. main detects it via errors.As to print starter UX.
type MissingPoolConfigError struct {
	Path string
}

func (e *MissingPoolConfigError) Error() string {
	return fmt.Sprintf("pool config not found at %s", e.Path)
}

func LoadPool() (*PoolConfig, error) {
	path := PoolConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, &MissingPoolConfigError{Path: path}
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	if err := CheckSchemaVersion(raw, path, PoolSchemaVersion); err != nil {
		return nil, err
	}

	cfg := PoolConfig{
		PortRangeStart: DefaultPortRangeStart,
		PortRangeEnd:   DefaultPortRangeEnd,
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", path, err)
	}
	if err := cfg.validate(path); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *PoolConfig) validate(path string) error {
	if c.PortRangeStart <= 0 {
		return fmt.Errorf("invalid %s:\n  - portRangeStart: must be a positive integer", path)
	}
	if c.PortRangeEnd <= 0 {
		return fmt.Errorf("invalid %s:\n  - portRangeEnd: must be a positive integer", path)
	}
	if c.PortRangeEnd < c.PortRangeStart {
		return fmt.Errorf(
			"invalid %s:\n  - <root>: portRangeEnd (%d) must be >= portRangeStart (%d)",
			path, c.PortRangeEnd, c.PortRangeStart,
		)
	}
	if c.PortRangeEnd < MinPort || c.PortRangeStart > MaxPort {
		return fmt.Errorf(
			"invalid %s:\n  - <root>: range %d-%d has no usable ports (port-pool only hands out %d-%d)",
			path, c.PortRangeStart, c.PortRangeEnd, MinPort, MaxPort,
		)
	}
	for i, p := range c.ExcludedPorts {
		if p <= 0 {
			return fmt.Errorf("invalid %s:\n  - excludedPorts.%d: must be a positive integer", path, i)
		}
	}
	return nil
}
