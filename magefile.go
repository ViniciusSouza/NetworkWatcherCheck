// +build mage

// Build a script to format and run tests of a Terraform module project
package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// The default target when the command executes `mage` in Cloud Shell
var Default = Full

// A build step that runs Clean, Format, Unit and Integration in sequence
func Full() {
	mg.Deps(UnitTest)
	mg.Deps(Integration)
}

// A build step that runs unit tests
func UnitTest() error {
	mg.Deps(Format)
	fmt.Println("Running Unit tests...")
	return sh.RunV("go", "test", "./...", "-run", "TestUN_", "-v")
}

// A build step that runs integration tests
func Integration() error {
	mg.Deps(Format)
	fmt.Println("Running integration tests...")
	return sh.RunV("go", "test", "./...", "-run", "TestIT_", "-v")
}

// A build step that formats Go code
func Format() error {
	fmt.Println("Formatting...")
	return sh.RunV("go", "fmt", "./...")
}
