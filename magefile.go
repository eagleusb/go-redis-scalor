//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

var Aliases = map[string]interface{} {
  "r":     StartRedis,
  "s":     StopRedis,
}

var (
	defaultRedisHost string = "0"
	defaultRedisPort string = "32000"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = StartRedis

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

func ListRedis() (err error) {
	fmt.Print("Listing testing Redis shards\n")
	err = sh.RunV("redis-cli", "-h", defaultRedisHost, "-p", defaultRedisPort, "cluster", "nodes")
	return
}
