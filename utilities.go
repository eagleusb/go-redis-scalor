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

// TODO
// func rebalanceRedis()  {
// redis-cli --cluster reshard <host>:<port> --cluster-from <node-id> --cluster-to <node-id> --cluster-slots <number of slots> --cluster-yes
// }

func (c *RedisClusterConf) wantedRedisConfArg(lookup string) bool {
	switch lookup {
	case
		"cluster_state",
		"cluster_slots_assigned",
		"cluster_slots_ok",
		"cluster_known_nodes",
		"cluster_size":
		return true
	}
	return false
}

func (c *RedisClusterConf) setRedisConfArg(arg []string) {
	switch arg[0] {
	case "cluster_state":
			c.State = arg[1]
	case "cluster_slots_assigned":
			c.SlotsAssigned = arg[1]
	case "cluster_slots_ok":
			c.SlotsOk = arg[1]
	case "cluster_known_nodes":
			c.NodesKnown = arg[1]
	case "cluster_size":
			c.Size = arg[1]
	}
}

func (c *RedisClusterConf) getRedisConfArg(arg string) (_value string) {
	// _value = ""
	switch arg {
	case "cluster_state":
			_value = c.State
	case "cluster_slots_assigned":
			_value = c.SlotsAssigned
	case "cluster_slots_ok":
			_value = c.SlotsOk
	case "cluster_known_nodes":
			_value = c.NodesKnown
	case "cluster_size":
			_value = c.Size
	}
	return
}
