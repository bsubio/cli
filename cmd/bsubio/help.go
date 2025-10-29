package main

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed static/*.md
var staticFiles embed.FS

func runHelpCommand(args []string) error {
	if len(args) == 0 {
		return runHelp(args)
	}

	command := args[0]
	filename := fmt.Sprintf("static/%s.md", command)

	// Try to read the file from embedded FS
	content, err := staticFiles.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("no help available for command: %s", command)
	}

	fmt.Print(strings.TrimSpace(string(content)))
	fmt.Println()

	return nil
}

func runQuickstart(args []string) error {
	content, err := staticFiles.ReadFile("static/quickstart.md")
	if err != nil {
		return fmt.Errorf("quickstart guide not available")
	}

	fmt.Print(strings.TrimSpace(string(content)))
	fmt.Println()
	return nil
}
