package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/elpdev/bubbleplate/internal/app"
	"github.com/elpdev/bubbleplate/internal/generator"
	"github.com/elpdev/bubbleplate/internal/generatorui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	root := flag.NewFlagSet("bubbleplate", flag.ExitOnError)
	showVersion := root.Bool("version", false, "print version information")
	root.Usage = printUsage
	if err := root.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	if *showVersion {
		fmt.Printf("bubbleplate %s (%s, %s)\n", version, commit, date)
		return
	}

	args := root.Args()
	if len(args) == 0 {
		runGeneratorUI()
		return
	}

	switch args[0] {
	case "demo":
		runDemo()
	case "new":
		runNew(args[1:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		printUsage()
		os.Exit(2)
	}
}

func runGeneratorUI() {
	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		printUsage()
		return
	}
	program := tea.NewProgram(generatorui.New())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "bubbleplate: %v\n", err)
		os.Exit(1)
	}
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func runDemo() {
	meta := app.BuildInfo{Version: version, Commit: commit, Date: date}
	program := tea.NewProgram(app.New(meta))
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "bubbleplate: %v\n", err)
		os.Exit(1)
	}
}

func runNew(args []string) {
	flags := flag.NewFlagSet("bubbleplate new", flag.ExitOnError)
	modulePath := flags.String("module", "", "Go module path for the generated project")
	outputDir := flags.String("output", "", "directory to create; defaults to app name")
	description := flags.String("description", "", "project description")
	displayName := flags.String("display-name", "", "human-readable app name")
	dockerImage := flags.String("docker-image", "", "Docker image shown in generated README")
	force := flags.Bool("force", false, "allow writing into a non-empty output directory")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage: bubbleplate new <name> --module <module-path> [options]\n\n")
		flags.PrintDefaults()
	}
	if len(args) == 0 {
		flags.Usage()
		os.Exit(2)
	}
	appName := args[0]
	if err := flags.Parse(args[1:]); err != nil {
		os.Exit(2)
	}
	if flags.NArg() != 0 {
		flags.Usage()
		os.Exit(2)
	}

	config := generator.NewConfig(appName)
	config.ModulePath = *modulePath
	config.Force = *force
	if *outputDir != "" {
		config.OutputDir = *outputDir
	}
	if *description != "" {
		config.Description = *description
	}
	if *displayName != "" {
		config.DisplayName = *displayName
	}
	if *dockerImage != "" {
		config.DockerImage = *dockerImage
	}

	result, err := generator.Generate(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bubbleplate new: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s (%d files)\n", result.OutputDir, result.Files)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", result.OutputDir)
	fmt.Println("  go mod tidy")
	fmt.Println("  go test ./...")
	fmt.Printf("  go run ./cmd/%s\n", config.BinaryName)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Bubbleplate generates opinionated Bubble Tea TUI projects.

Usage:
  bubbleplate
  bubbleplate new <name> --module <module-path> [options]
  bubbleplate demo
  bubbleplate --version
  bubbleplate --help

Examples:
  bubbleplate
  bubbleplate new myapp --module github.com/acme/myapp
  bubbleplate new myapp --module github.com/acme/myapp --output ../myapp
  bubbleplate demo

`)
}
