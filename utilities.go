package main

import (
	"fmt"
	"os/exec"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO
// func newNode() *RedisNode { return }

func slotsCount(slotStart, slotEnd int) int {
 return ((slotEnd + 1) - slotStart)
}

func slotsPercentage(slots int) int {
	return int((float64(slots) / float64(16384)) * float64(100))
}

func isRedisCli() (path string) {
	path, err := exec.LookPath("redis-cli")
	checkErr(err)
	fmt.Println("Found redis-cli at", path)
	return
}

func execRedisCli(args string) {
	stdout, err := exec.Command(isRedisCli(), args).Output()
	checkErr(err)
	fmt.Printf("%v\n", string(stdout))
}
