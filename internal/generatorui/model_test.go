package generatorui

import "testing"

func TestConfigFromValuesUsesDefaults(t *testing.T) {
	m := New()
	m.values[fieldAppName] = "hackernews"
	m.values[fieldModulePath] = "github.com/elpdev/hackernews"

	config, err := m.configFromValues()
	if err != nil {
		t.Fatalf("configFromValues failed: %v", err)
	}

	if config.AppName != "hackernews" {
		t.Fatalf("AppName = %q, want hackernews", config.AppName)
	}
	if config.BinaryName != "hackernews" {
		t.Fatalf("BinaryName = %q, want hackernews", config.BinaryName)
	}
	if config.OutputDir != "hackernews" {
		t.Fatalf("OutputDir = %q, want hackernews", config.OutputDir)
	}
	if config.DisplayName != "Hackernews" {
		t.Fatalf("DisplayName = %q, want Hackernews", config.DisplayName)
	}
	if config.ModulePath != "github.com/elpdev/hackernews" {
		t.Fatalf("ModulePath = %q, want github.com/elpdev/hackernews", config.ModulePath)
	}
}

func TestConfigFromValuesUsesOverrides(t *testing.T) {
	m := New()
	m.values[fieldAppName] = "hackernews"
	m.values[fieldModulePath] = "github.com/elpdev/hackernews"
	m.values[fieldOutputDir] = "../hn"
	m.values[fieldDisplayName] = "Hacker News"
	m.values[fieldDescription] = "A terminal Hacker News reader"
	m.values[fieldDockerImage] = "ghcr.io/elpdev/hackernews"
	m.force = true

	config, err := m.configFromValues()
	if err != nil {
		t.Fatalf("configFromValues failed: %v", err)
	}

	if config.OutputDir != "../hn" {
		t.Fatalf("OutputDir = %q, want ../hn", config.OutputDir)
	}
	if config.DisplayName != "Hacker News" {
		t.Fatalf("DisplayName = %q, want Hacker News", config.DisplayName)
	}
	if config.Description != "A terminal Hacker News reader" {
		t.Fatalf("Description = %q, want custom description", config.Description)
	}
	if config.DockerImage != "ghcr.io/elpdev/hackernews" {
		t.Fatalf("DockerImage = %q, want ghcr.io/elpdev/hackernews", config.DockerImage)
	}
	if !config.Force {
		t.Fatal("Force = false, want true")
	}
}

func TestConfigFromValuesRequiresAppNameAndModule(t *testing.T) {
	m := New()
	if _, err := m.configFromValues(); err == nil {
		t.Fatal("expected missing app name error")
	}

	m.values[fieldAppName] = "hackernews"
	if _, err := m.configFromValues(); err == nil {
		t.Fatal("expected missing module path error")
	}
}
