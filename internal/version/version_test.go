package version

import (
	"strings"
	"testing"
)

func TestFullVersionContainsFields(t *testing.T) {
	out := FullVersion()

	expected := []string{
		"Flint Version:",
		"Commit:",
		"Build Time:",
		"Go Version:",
		"Platform:",
	}

	for _, field := range expected {
		if !strings.Contains(out, field) {
			t.Fatalf("expected output to contain %q, got:\n%s", field, out)
		}
	}
}

func TestDefaultsAreUsed(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}
	if GitCommit == "" {
		t.Fatal("GitCommit must not be empty")
	}
	if BuildTime == "" {
		t.Fatal("BuildTime must not be empty")
	}
	if GoVersion == "" {
		t.Fatal("GoVersion must not be empty")
	}
	if OS == "" {
		t.Fatal("OS must not be empty")
	}
	if Arch == "" {
		t.Fatal("Arch must not be empty")
	}
}

func TestFullVersionFormatting(t *testing.T) {
	out := FullVersion()

	lines := strings.Split(out, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d:\n%s", len(lines), out)
	}

	if !strings.HasPrefix(lines[0], "Flint Version:") {
		t.Fatal("line 1 should begin with 'Flint Version:'")
	}
	if !strings.HasPrefix(lines[4], "Platform:") {
		t.Fatal("last line should begin with 'Platform:'")
	}
}
