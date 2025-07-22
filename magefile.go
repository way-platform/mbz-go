//go:build mage

package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build runs a full CI build.
func Build() {
	mg.Deps(Download)
	mg.Deps(Generate)
	mg.Deps(Lint, Test)
	mg.Deps(Tidy)
	mg.Deps(Diff)
}

// Lint runs the Go linter.
func Lint() error {
	return sh.RunV("go", "tool", "golangci-lint", "run")
}

// Test runs the Go tests.
func Test() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "test", "-v", "-cover", "./...").Run()
	})
}

// Download downloads the Go dependencies.
func Download() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "download").Run()
	})
}

// Tidy tidies the Go mod files.
func Tidy() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "tidy", "-v").Run()
	})
}

// Diff checks for git diffs.
func Diff() error {
	return sh.RunV("git", "diff", "--exit-code")
}

// Generate runs all code generators.
func Generate() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "generate", "-v", "./...").Run()
	})
}

// MBZ builds the mbz CLI.
func MBZ() error {
	return cmd("cmd/mbz", "go", "install", ".").Run()
}

// VHS records the CLI GIF using VHS.
func VHS() error {
	mg.Deps(MBZ)
	return cmd("docs", "go", "tool", "vhs", "cli.tape").Run()
}

func forEachGoMod(f func(dir string) error) error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "go.mod" {
			return nil
		}
		return f(filepath.Dir(path))
	})
}

func cmd(dir string, command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
