//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Runs go mod download and then installs the binary.
func Lint() error {
	return sh.Run("golangci-lint", "-v", "run")
}

func Fmt() error {
	return sh.Run("gofumpt", "-w", "-l", ".")
}

func Hello() error {
	fmt.Println("Hello")
	return nil
}
