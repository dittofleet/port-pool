package cmd

import (
	"fmt"
	"os"
	"strings"

	clikit "github.com/dittofleet/go-cli-kit"
	"github.com/dittofleet/go-cli-kit/uninstall"
	"github.com/dittofleet/go-cli-kit/xdg"
	"github.com/dittofleet/port-pool/internal/state"
)

const uninstallUsage = "usage: port-pool uninstall [--yes]"

// Uninstall removes the port-pool binary, global config directory, and
// data directory.
func Uninstall(args []string, a clikit.App) error {
	args, yes := extractBoolFlag(args, "yes")
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			return fmt.Errorf("unknown flag: %s\n%s", arg, uninstallUsage)
		}
	}
	if len(args) > 0 {
		return fmt.Errorf("unexpected arguments: %v\n%s", args, uninstallUsage)
	}

	var allocations string
	if s, err := state.Load(); err == nil {
		switch n := len(s.Allocations); n {
		case 0:
			allocations = "no active allocations"
		case 1:
			allocations = "1 active allocation"
		default:
			allocations = fmt.Sprintf("%d active allocations", n)
		}
	}

	return uninstall.Run(a, yes, uninstall.Plan{
		Items: []uninstall.Item{
			{Label: "State", Path: xdg.DataDir(a.Name), Note: allocations, Remove: os.RemoveAll},
			{Label: "Config", Path: xdg.ConfigDir(a.Name), Remove: os.RemoveAll},
		},
		Notice: "Note: .env files in project directories are NOT touched.",
	})
}
