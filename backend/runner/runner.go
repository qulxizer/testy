package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"testy/models"
)

const goModContent = `module student

go 1.20

require github.com/01-edu/z01 v0.1.0
`

func RunSubmission(ctx context.Context, req models.SubmissionRequest) (string, bool, error) {
	tmpDir, err := os.MkdirTemp("", "sub_*")
	if err != nil {
		return "", false, err
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		return "", false, err
	}

	exerciseDir := filepath.Join(tmpDir, req.Exercise)
	if err := os.MkdirAll(exerciseDir, 0755); err != nil {
		return "", false, err
	}

	if len(req.Files) > 0 {
		for fname, content := range req.Files {
			if req.Filename != "main.go" && fname == "main.go" {
				continue
			}
			_ = os.WriteFile(filepath.Join(exerciseDir, fname), []byte(content), 0644)
			_ = os.WriteFile(filepath.Join(tmpDir, fname), []byte(content), 0644)
		}
	} else {
		codeBytes := []byte(req.Code)
		_ = os.WriteFile(filepath.Join(exerciseDir, req.Filename), codeBytes, 0644)
		_ = os.WriteFile(filepath.Join(tmpDir, req.Filename), codeBytes, 0644)
	}

	testCmd := `cp -r /go-tests ~ && cd ~/go-tests && if command -v "${EXERCISE}_test" >/dev/null 2>&1; then "${EXERCISE}_test"; else go run "./tests/${EXERCISE}_test"; fi`

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1.0",
		"--pids-limit", "100",
		"-e", "EXERCISE="+req.Exercise,
		"-v", tmpDir+":/root/piscine-go:ro",
		"--entrypoint", "sh",
		"ghcr.io/01-edu/test-go:latest",
		"-c", testCmd,
	)

	out, err := cmd.CombinedOutput()
	passed := err == nil && cmd.ProcessState != nil && cmd.ProcessState.Success()

	outputStr := string(out)
	if passed && outputStr == "" {
		outputStr = "✓ All tests passed successfully!"
	}

	return outputStr, passed, nil
}
