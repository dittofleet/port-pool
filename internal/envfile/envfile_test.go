package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCRLFFileIsNotDriftAndKeepsItsLineEndings(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("PORT=3100\r\nSECRET=abc\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	drifted, err := Diff(path, map[string]string{"PORT": "3100"})
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Fatalf("expected no drift, got %v", drifted)
	}

	if err := Apply(path, map[string]string{"PORT": "3200", "API": "3201"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "PORT=3200\r\nSECRET=abc\r\nAPI=3201\r\n"
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBOMDoesNotHideTheFirstKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("\xef\xbb\xbfPORT=3100\nSECRET=abc\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	drifted, err := Diff(path, map[string]string{"PORT": "3100"})
	if err != nil {
		t.Fatal(err)
	}
	if len(drifted) != 0 {
		t.Fatalf("expected no drift, got %v", drifted)
	}

	if err := Apply(path, map[string]string{"PORT": "3200"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "PORT=3200\nSECRET=abc\n"; string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplyWritesLFToANewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := Apply(path, map[string]string{"PORT": "3100"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "PORT=3100\n"; string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
