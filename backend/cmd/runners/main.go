package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	repoURL       = "https://github.com/01-edu/public"
	repoDir       = "./tmp/01edu"
	runnersJSON   = "../frontend/data/runners.json"
	stubsJSON     = "../frontend/data/stubs.json"
	subjectsJSON  = "../frontend/data/subjects.json"
	auditsJSON    = "../frontend/data/audits.json"
	exercisesJSON = "../frontend/data/exercises.json"
	auditAPIURL   = "https://learn.reboot01.com/api/object/bahrain"
)

type Exercise struct {
	ID               int      `json:"id"`
	Key              string   `json:"key"`
	Name             string   `json:"name"`
	Path             string   `json:"path"`
	Difficulty       float64  `json:"difficulty"`
	Level            int      `json:"level"`
	ExpectedFiles    []string `json:"expectedFiles"`
	AllowedFunctions []string `json:"allowedFunctions"`
	TestImage        string   `json:"testImage,omitempty"`
	Subject          string   `json:"subject"`
	BaseXP           int      `json:"baseXp"`
}

type APIObject struct {
	ID       int                   `json:"id"`
	Name     string                `json:"name"`
	Type     string                `json:"type"`
	Key      string                `json:"key"`
	Path     string                `json:"path"`
	Attrs    Attrs                 `json:"attrs"`
	Children map[string]*APIObject `json:"children"`
}

type Attrs struct {
	Subject          string       `json:"subject"`
	Difficulty       float64      `json:"difficulty"`
	Level            int          `json:"level"`
	BaseXP           int          `json:"baseXp"`
	ExpectedFiles    []string     `json:"expectedFiles"`
	AllowedFunctions []string     `json:"allowedFunctions"`
	Validations      []Validation `json:"validations"`
}

type Validation struct {
	TestImage string `json:"testImage"`
}

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

func parseExercises(obj *APIObject, list *[]Exercise) {
	if obj == nil {
		return
	}

	if obj.Type == "exercise" {
		testImg := ""
		if len(obj.Attrs.Validations) > 0 {
			testImg = obj.Attrs.Validations[0].TestImage
		}

		*list = append(*list, Exercise{
			ID:               obj.ID,
			Key:              obj.Key,
			Name:             obj.Name,
			Path:             obj.Path,
			Difficulty:       obj.Attrs.Difficulty,
			Level:            obj.Attrs.Level,
			ExpectedFiles:    obj.Attrs.ExpectedFiles,
			AllowedFunctions: obj.Attrs.AllowedFunctions,
			TestImage:        testImg,
			Subject:          obj.Attrs.Subject,
			BaseXP:           obj.Attrs.BaseXP,
		})
	}

	for _, child := range obj.Children {
		parseExercises(child, list)
	}
}

func fetchAndExtractExercises() error {
	fmt.Println("Fetching Bahrain API object...")
	resp, err := http.Get(auditAPIURL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var root APIObject
	if err := json.Unmarshal(body, &root); err != nil {
		return fmt.Errorf("failed to unmarshal root: %w", err)
	}

	var exercises []Exercise
	parseExercises(&root, &exercises)

	if err := os.MkdirAll(filepath.Dir(exercisesJSON), 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	outData, err := json.MarshalIndent(exercises, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal exercises: %w", err)
	}

	if err := os.WriteFile(exercisesJSON, outData, 0644); err != nil {
		return fmt.Errorf("failed to write exercises.json: %w", err)
	}

	fmt.Printf("✓ Extracted %d exercises to %s\n", len(exercises), exercisesJSON)
	return nil
}

func writeJSONFile(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func generateData() error {
	subjectsDir := filepath.Join(repoDir, "subjects")
	entries, err := os.ReadDir(subjectsDir)
	if err != nil {
		return fmt.Errorf("failed to read subjects dir: %w", err)
	}

	runners := make(map[string]string)
	stubs := make(map[string]string)
	subjects := make(map[string]string)
	audits := make(map[string]string)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		folder := filepath.Join(subjectsDir, name)

		// 1. Runner main.go
		if content, err := os.ReadFile(filepath.Join(folder, "main.go")); err == nil {
			runners[name] = string(content)
		}

		// 2. Subject README.md
		readmePath := filepath.Join(folder, "README.md")
		if content, err := os.ReadFile(readmePath); err == nil {
			subjects[name] = string(content)
			if stub := extractStub(readmePath); stub != "" {
				stubs[name] = stub
			}
		}

		// 3. Audit README.md (if present in /audit/README.md)
		auditPath := filepath.Join(folder, "audit", "README.md")
		if content, err := os.ReadFile(auditPath); err == nil {
			audits[name] = string(content)
		}
	}

	if err := writeJSONFile(runnersJSON, runners); err != nil {
		return fmt.Errorf("failed writing runners: %w", err)
	}
	if err := writeJSONFile(stubsJSON, stubs); err != nil {
		return fmt.Errorf("failed writing stubs: %w", err)
	}
	if err := writeJSONFile(subjectsJSON, subjects); err != nil {
		return fmt.Errorf("failed writing subjects markdown: %w", err)
	}
	if err := writeJSONFile(auditsJSON, audits); err != nil {
		return fmt.Errorf("failed writing audits markdown: %w", err)
	}

	if err := fetchAndExtractExercises(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	fmt.Printf("✓ Extracted %d runners to %s\n", len(runners), runnersJSON)
	fmt.Printf("✓ Extracted %d function stubs to %s\n", len(stubs), stubsJSON)
	fmt.Printf("✓ Extracted %d subject markdowns to %s\n", len(subjects), subjectsJSON)
	fmt.Printf("✓ Extracted %d audit markdowns to %s\n", len(audits), auditsJSON)
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
	fmt.Println("  clone   Clone repo, extract markdowns/runners/stubs, and fetch exercise metadata")
	fmt.Println("  pull    Pull latest changes and rebuild all JSON files")
	fmt.Println("  sync    Rebuild all JSON files from existing clone and API")
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
