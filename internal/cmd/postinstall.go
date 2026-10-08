package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/postinstall"
	"github.com/dittofleet/port-pool/internal/config"
)

// Postinstall is the install script's first-time setup: the starter
// config, unless there is one already.
func Postinstall(a clikit.App) error {
	return postinstall.Run(a, func() error {
		path := config.PoolConfigPath()
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(config.StarterPoolConfig+"\n"), 0o644); err != nil {
			return err
		}
		fmt.Printf("Created starter config at %s\n", path)
		return nil
	})
}
