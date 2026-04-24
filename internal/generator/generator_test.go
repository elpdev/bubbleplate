package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCreatesNamedProject(t *testing.T) {
	output := filepath.Join(t.TempDir(), "myapp")
	config := NewConfig("myapp")
	config.ModulePath = "github.com/acme/myapp"
	config.OutputDir = output
	config.DisplayName = "My App"
	config.Description = "My generated TUI"
	config.DockerImage = "ghcr.io/acme/myapp"

	result, err := Generate(config)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if result.Files == 0 {
		t.Fatal("expected generated files")
	}

	assertFileContains(t, output, "go.mod", "module github.com/acme/myapp")
	assertFileContains(t, output, filepath.Join("cmd", "myapp", "main.go"), "github.com/acme/myapp/internal/app")
	assertFileContains(t, output, filepath.Join("cmd", "myapp", "main.go"), "myapp %s")
	assertFileContains(t, output, ".goreleaser.yaml", "project_name: myapp")
	assertFileContains(t, output, "Dockerfile", "/out/myapp")
	assertFileContains(t, output, filepath.Join(".github", "workflows", "publish.yml"), "ghcr.io/${{ github.repository_owner }}/myapp")
	assertFileContains(t, output, "README.md", "# My App")
	assertFileContains(t, output, "README.md", "My generated TUI")
	assertFileContains(t, output, "README.md", "ghcr.io/acme/myapp:latest")
	assertFileWritable(t, output, "go.mod")
}

func TestGenerateRefusesNonEmptyOutputWithoutForce(t *testing.T) {
	output := t.TempDir()
	if err := os.WriteFile(filepath.Join(output, "existing.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	config := NewConfig("myapp")
	config.ModulePath = "github.com/acme/myapp"
	config.OutputDir = output

	_, err := Generate(config)
	if err == nil {
		t.Fatal("expected non-empty output error")
	}
}

func TestGeneratedProjectBuildsAndTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping generated project integration test in short mode")
	}

	output := filepath.Join(t.TempDir(), "sampleapp")
	config := NewConfig("sampleapp")
	config.ModulePath = "github.com/acme/sampleapp"
	config.OutputDir = output
	config.DisplayName = "Sample App"

	if _, err := Generate(config); err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	run(t, output, "go", "test", "./...")
	run(t, output, "go", "build", "./cmd/sampleapp")
}

func assertFileContains(t *testing.T, root, name, want string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	if !strings.Contains(string(content), want) {
		t.Fatalf("expected %s to contain %q", name, want)
	}
}

func assertFileWritable(t *testing.T, root, name string) {
	t.Helper()
	info, err := os.Stat(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("stat %s: %v", name, err)
	}
	if info.Mode().Perm()&0o200 == 0 {
		t.Fatalf("expected %s to be owner-writable, got mode %v", name, info.Mode().Perm())
	}
}

func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, string(out))
	}
}
