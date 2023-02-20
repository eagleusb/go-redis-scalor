//go:build mage
// +build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/sh"
)

var Aliases = map[string]interface{} {
  "r":     StartRedis,
  "s":     StopRedis,
}

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = StartRedis

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

func StartRedis() (err error) {
	fmt.Print("Launching testing Redis cluster")
	err = sh.RunV("docker-compose", "up", "--detach", "--remove-orphans")
	return
}

func StopRedis() (err error) {
	fmt.Print("Stopping testing Redis cluster")
	err = sh.RunV("docker-compose", "down", "--remove-orphans")
	return
}
