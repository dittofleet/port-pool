package config

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func loadProjectWithEnvFile(t *testing.T, envPath string) error {
	t.Helper()
	dir := t.TempDir()
	data, err := json.Marshal(map[string]any{
		"schemaVersion": 1,
		"portNames":     []string{"web"},
		"envFiles":      map[string]any{envPath: map[string]string{"PORT": "${web}"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ProjectConfigPath(dir), data, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = LoadProject(dir)
	return err
}

func TestLoadProjectAcceptsPathsInsideTheProject(t *testing.T) {
	for _, envPath := range []string{".env", "apps/web/.env.local"} {
		if err := loadProjectWithEnvFile(t, envPath); err != nil {
			t.Errorf("%s: %v", envPath, err)
		}
	}
}

func assertOutsideProject(t *testing.T, envPath string) {
	t.Helper()
	err := loadProjectWithEnvFile(t, envPath)
	if err == nil || !strings.Contains(err.Error(), "must be a relative path") {
		t.Errorf("%s: got %v, want a relative path error", envPath, err)
	}
}

func TestLoadProjectRejectsAbsolutePaths(t *testing.T) {
	assertOutsideProject(t, "/etc/app/.env")
}
