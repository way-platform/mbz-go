//go:build mage

package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build runs a full CI build.
func Build() {
	mg.Deps(Generate, Lint, Test, Tidy, Diff)
}

// Lint runs the Go linter.
func Lint() error {
	return sh.RunV("go", "tool", "golangci-lint", "run")
}

// Test runs the Go tests.
func Test() error {
	mg.Deps(Generate)
	return sh.RunV("go", "test", "-v", "-cover", "./...")
}

// Tidy tidies the Go mod files.
func Tidy() error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "go.mod" {
			return nil
		}
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		defer os.Chdir(currentDir)
		if err := os.Chdir(filepath.Dir(path)); err != nil {
			return err
		}
		return sh.RunV("go", "mod", "tidy", "-v")
	})
}

// Diff checks for git diffs.
func Diff() error {
	return sh.RunV("git", "diff", "--exit-code")
}

// Generate runs all code generators.
func Generate() error {
	return sh.RunV("go", "generate", "-v", "./...")
}
