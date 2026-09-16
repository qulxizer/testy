package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"testy/models"
)

const goModContent = `module student

go 1.26

require github.com/01-edu/z01 v0.1.0
`

const testCommand = `
cp -r /go-tests ~ &&
cd ~/go-tests &&
if command -v "${EXERCISE}_test" >/dev/null 2>&1; then
	"${EXERCISE}_test"
else
	go run "./tests/${EXERCISE}_test"
fi
`

func RunSubmission(ctx context.Context, req models.SubmissionRequest) (string, bool, error) {
	tmpDir, err := createSubmissionDir()
	if err != nil {
		return "", false, err
	}
	fmt.Println("tmpDir:", tmpDir)

	// defer os.RemoveAll(tmpDir)

	if err := writeSubmission(tmpDir, req); err != nil {
		return "", false, err
	}

	output, err := runTests(ctx, tmpDir, req.Exercise, req.TestImage)
	if err != nil {
		return string(output), false, nil
	}

	if len(output) == 0 {
		output = []byte("✓ All tests passed successfully!")
	}

	return string(output), true, nil
}

func createSubmissionDir() (string, error) {
	return os.MkdirTemp("/submissions", "sub_*")
}

func writeSubmission(dir string, req models.SubmissionRequest) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create submission directory: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(dir, "go.mod"),
		[]byte(goModContent),
		0644,
	); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	if len(req.Files) > 0 {
		if err := writeFiles(dir, req); err != nil {
			return fmt.Errorf("write submission files: %w", err)
		}
	} else {
		if err := writeCode(dir, req); err != nil {
			return fmt.Errorf("write submission code: %w", err)
		}
	}

	return nil
}

func writeFiles(dir string, req models.SubmissionRequest) error {
	for filename, content := range req.Files {
		filename = filepath.FromSlash(filename)

		if filename == "main.go" {
			continue
		}

		path := filepath.Join(dir, filename)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func writeCode(dir string, req models.SubmissionRequest) error {
	path := filepath.Join(dir, req.Filename)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(req.Code), 0644); err != nil {
		return fmt.Errorf("write %s: %w", req.Filename, err)
	}

	return nil
}

func runTests(ctx context.Context, tmpDir string, exercise string, testImage string) ([]byte, error) {
	submission := filepath.Base(tmpDir)

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1.0",
		"--pids-limit", "100",

		"-e", "EXERCISE="+exercise,
		"-e", "SUBMISSION="+submission,

		"-v", "testy_submissions:/submissions:ro",

		"--entrypoint", "sh",
		testImage,

		"-c", `
			rm -rf /root/piscine-go
			cp -r "/submissions/$SUBMISSION" /root/piscine-go
			`+testCommand,
	)

	return cmd.CombinedOutput()
}
