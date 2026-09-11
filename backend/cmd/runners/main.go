package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	repoURL     = "https://github.com/01-edu/public"
	repoDir     = "./tmp/01edu"
	runnersJSON = "../frontend/app/data/runners.json"
	stubsJSON   = "../frontend/app/data/stubs.json"
)

// Grabs Go code block declaring the function (and any custom types before it)
var stubRegex = regexp.MustCompile("(?s)```(?:go)?\\s*\\n((?:type\\s+[\\s\\S]*?\\n\\s*)?func\\s+[A-Z][^`]*?)```")

func extractStub(readmePath string) string {
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return ""
	}

	matches := stubRegex.FindSubmatch(content)
	if len(matches) < 2 {
		return ""
	}

	raw := strings.TrimSpace(string(matches[1]))
	if raw == "" {
		return ""
	}

	// Add empty brackets if the subject only provided the signature
	if !strings.Contains(raw, "{") {
		raw += " {\n\t\n}"
	}

	return fmt.Sprintf("package piscine\n\n%s\n", raw)
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateData() error {
	subjectsDir := filepath.Join(repoDir, "subjects")
	entries, err := os.ReadDir(subjectsDir)
	if err != nil {
		return fmt.Errorf("failed to read subjects dir: %w", err)
	}

	runners := make(map[string]string)
	stubs := make(map[string]string)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		folder := filepath.Join(subjectsDir, name)

		// 1. Grab runner main.go if present
		mainFile := filepath.Join(folder, "main.go")
		if content, err := os.ReadFile(mainFile); err == nil {
			runners[name] = string(content)
		}

		// 2. Extract function prototype from README.md
		readmeFile := filepath.Join(folder, "README.md")
		if stub := extractStub(readmeFile); stub != "" {
			stubs[name] = stub
		}
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(runnersJSON), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write runners.json
	runnersData, err := json.MarshalIndent(runners, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal runners: %w", err)
	}
	if err := os.WriteFile(runnersJSON, runnersData, 0644); err != nil {
		return fmt.Errorf("failed to write runners json: %w", err)
	}

	// Write stubs.json
	stubsData, err := json.MarshalIndent(stubs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stubs: %w", err)
	}
	if err := os.WriteFile(stubsJSON, stubsData, 0644); err != nil {
		return fmt.Errorf("failed to write stubs json: %w", err)
	}

	fmt.Printf("✓ Extracted %d runners to %s\n", len(runners), runnersJSON)
	fmt.Printf("✓ Extracted %d function stubs to %s\n", len(stubs), stubsJSON)
	return nil
}

func cloneRepo() {
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
		fmt.Printf("Repo already exists at %s\n", repoDir)
		_ = generateData()
		return
	}

	if _, err := os.Stat(repoDir); err == nil {
		fmt.Fprintf(os.Stderr, "Cannot clone: %s exists but is not a git repo\n", repoDir)
		os.Exit(1)
	}

	_ = os.MkdirAll(filepath.Dir(repoDir), 0755)
	fmt.Println("Cloning public repo (shallow)...")
	if err := runCmd("git", "clone", "--depth", "1", repoURL, repoDir); err != nil {
		fmt.Fprintf(os.Stderr, "Clone failed: %v\n", err)
		os.Exit(1)
	}

	if err := generateData(); err != nil {
		fmt.Fprintf(os.Stderr, "Generation failed: %v\n", err)
	}
}

func pullRepo() {
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err != nil {
		fmt.Fprintf(os.Stderr, "Repo not found at %s. Run 'clone' first.\n", repoDir)
		os.Exit(1)
	}

	fmt.Println("Pulling latest changes...")
	if err := runCmd("git", "-C", repoDir, "pull", "--ff-only"); err != nil {
		fmt.Fprintf(os.Stderr, "Pull failed: %v\n", err)
		os.Exit(1)
	}

	if err := generateData(); err != nil {
		fmt.Fprintf(os.Stderr, "Generation failed: %v\n", err)
	}
}

func printUsage() {
	fmt.Println("Usage: go run main.go <clone|pull|sync>")
	fmt.Println()
	fmt.Println("  clone  Clone repo into local path and build JSON files")
	fmt.Println("  pull   Pull latest changes and rebuild JSON files")
	fmt.Println("  sync   Re-extract runners and stubs from existing local repo")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "clone":
		cloneRepo()
	case "pull":
		pullRepo()
	case "sync":
		if err := generateData(); err != nil {
			fmt.Fprintf(os.Stderr, "Sync failed: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}
