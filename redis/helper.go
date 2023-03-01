package rediscluster

import (
	"fmt"
	"os/exec"

	"github.com/eagleusb/go-redis-scalor/utilities"
)

func slotsCount(slotStart, slotEnd int) int {
	return ((slotEnd + 1) - slotStart)
}

func slotsPercentage(slots int) int {
	return int((float64(slots) / float64(16384)) * float64(100))
}

func isRedisCli() (path string) {
	path, err := exec.LookPath("redis-cli")
	utilities.CheckErr(err)
	fmt.Println("Found redis-cli at", path)
	return
}

func execRedisCli(args []string) string {
	cmd := exec.Command(isRedisCli(), args...)
	fmt.Println("Execution of command:", cmd)
	stdout, err := cmd.Output()
	utilities.CheckErr(err)
	return string(stdout)
}
