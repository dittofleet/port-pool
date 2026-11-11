package config

import "testing"

func TestLoadProjectRejectsPathsWindowsTreatsAsOutsideTheProject(t *testing.T) {
	for _, envPath := range []string{`C:\app\.env`, `C:.env`, `\app\.env`, `NUL`} {
		assertOutsideProject(t, envPath)
	}
}
