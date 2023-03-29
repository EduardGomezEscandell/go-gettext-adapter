//go:build tools
// +build tools

package main

// Copy-pasted from github.com/ubuntu/adsys

import (
	"os"

	"example.com/sample/internal/generators"
)

func main() {
	if generators.InstallOnlyMode() {
		os.Exit(1)
	}
}
