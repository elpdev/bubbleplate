package generator

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed testdata/starter/** testdata/starter/.dockerignore testdata/starter/.goreleaser.yaml testdata/starter/.github/workflows/*
var starterFS embed.FS

type Config struct {
	AppName     string
	DisplayName string
	BinaryName  string
	ModulePath  string
	Description string
	OutputDir   string
	DockerImage string
	Force       bool
}

type Result struct {
	OutputDir string
	Files     int
}

var appNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func NewConfig(appName string) Config {
	binaryName := strings.ToLower(strings.TrimSpace(appName))
	return Config{
		AppName:     strings.TrimSpace(appName),
		DisplayName: titleName(appName),
		BinaryName:  binaryName,
		Description: "A Bubble Tea terminal UI",
		OutputDir:   binaryName,
		DockerImage: "ghcr.io/<owner>/" + binaryName,
	}
}

func Generate(config Config) (Result, error) {
	config = normalize(config)
	if err := validate(config); err != nil {
		return Result{}, err
	}

	if err := ensureOutputDir(config.OutputDir, config.Force); err != nil {
		return Result{}, err
	}

	result := Result{OutputDir: config.OutputDir}
	err := fs.WalkDir(starterFS, "testdata/starter", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		rel, err := filepath.Rel("testdata/starter", path)
		if err != nil {
			return err
		}
		rel = strings.ReplaceAll(rel, "BINARY_NAME", config.BinaryName)
		rel = strings.TrimSuffix(rel, ".tmpl")
		destination := filepath.Join(config.OutputDir, filepath.FromSlash(rel))

		content, err := starterFS.ReadFile(path)
		if err != nil {
			return err
		}
		rendered, err := render(path, string(content), config)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(destination, rendered, fileMode(entry)); err != nil {
			return err
		}
		result.Files++
		return nil
	})
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func normalize(config Config) Config {
	config.AppName = strings.TrimSpace(config.AppName)
	if config.BinaryName == "" {
		config.BinaryName = strings.ToLower(config.AppName)
	}
	if config.DisplayName == "" {
		config.DisplayName = titleName(config.AppName)
	}
	if config.Description == "" {
		config.Description = "A Bubble Tea terminal UI"
	}
	if config.OutputDir == "" {
		config.OutputDir = config.BinaryName
	}
	if config.DockerImage == "" {
		config.DockerImage = "ghcr.io/<owner>/" + config.BinaryName
	}
	return config
}

func validate(config Config) error {
	if config.AppName == "" {
		return errors.New("app name is required")
	}
	if !appNamePattern.MatchString(config.AppName) {
		return fmt.Errorf("app name %q must contain only letters, numbers, hyphens, or underscores, and start with a letter or number", config.AppName)
	}
	if config.ModulePath == "" {
		return errors.New("module path is required; pass --module")
	}
	if strings.ContainsAny(config.BinaryName, `/\`) {
		return fmt.Errorf("binary name %q must not contain path separators", config.BinaryName)
	}
	return nil
}

func ensureOutputDir(path string, force bool) error {
	entries, err := os.ReadDir(path)
	if err == nil {
		if len(entries) > 0 && !force {
			return fmt.Errorf("output directory %q is not empty; pass --force to overwrite", path)
		}
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func render(path, content string, config Config) ([]byte, error) {
	tmpl, err := template.New(filepath.Base(path)).Delims("[[", "]]").Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", path, err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, config); err != nil {
		return nil, fmt.Errorf("render template %s: %w", path, err)
	}
	return out.Bytes(), nil
}

func fileMode(entry fs.DirEntry) os.FileMode {
	info, err := entry.Info()
	if err != nil {
		return 0o644
	}
	mode := info.Mode().Perm()
	if mode == 0 {
		return 0o644
	}
	return mode | 0o200
}

func titleName(name string) string {
	parts := strings.FieldsFunc(strings.TrimSpace(name), func(r rune) bool { return r == '-' || r == '_' || r == ' ' })
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}
