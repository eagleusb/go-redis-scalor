package main

import (
	"fmt"
	"os/exec"

	"github.com/redis/go-redis/v9"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
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

func execRedisCli(args []string) string {
	cmd := exec.Command(isRedisCli(), args...)
	fmt.Println("Execution of command:", cmd)
	stdout, err := cmd.Output()
	checkErr(err)
	return string(stdout)
}

func rebalanceRedisShard(master string, slots int, fromId string, toId string) {
	cmd := []string{"--cluster", "reshard", master, "--cluster-from", fromId, "--cluster-to", toId, "--cluster-slots", "100", "--cluster-yes"}
	execRedisCli(cmd)
}

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

func (r *RedisNodes) setRedisClusterNodes(client *redis.Client) {
	clusterSlots, _ := client.ClusterSlots(ctx).Result()

	for i, slot := range clusterSlots {
		_nodeId := clusterSlots[i].Nodes[0].ID
		_nodeIp := clusterSlots[i].Nodes[0].Addr

		clusterNode := &RedisNode{
			Id: _nodeId,
			Ip: _nodeIp,
		}

		if _, ok := r.ClusterNodes[_nodeId]; ok {
			clusterNode.Slots = append(r.ClusterNodes[_nodeId].Slots, fmt.Sprint(slot.Start, "-", slot.End))
			clusterNode.SlotsCount = r.ClusterNodes[_nodeId].SlotsCount + slotsCount(slot.Start, slot.End)
			clusterNode.SlotsPercentage = r.ClusterNodes[_nodeId].SlotsPercentage + slotsPercentage(slotsCount(slot.Start, slot.End))
		} else {
			clusterNode.Slots = make([]string, 1)
			clusterNode.Slots[0] = fmt.Sprint(slot.Start, "-", slot.End)
			clusterNode.SlotsCount = slotsCount(slot.Start, slot.End)
			clusterNode.SlotsPercentage = slotsPercentage(clusterNode.SlotsCount)
		}

		r.ClusterNodes[_nodeId] = *clusterNode
	}
}

func (r *RedisNodes) getRedisClusterNodes(client *redis.Client) {
		for _, node := range r.ClusterNodes {
		fmt.Printf(
			"ID: %s slots: %s slotsRanges: %d slotsCount: %d slotsPercentage: %d IP: %s\n",
			node.Id,
			node.Slots,
			len(node.Slots),
			node.SlotsCount,
			node.SlotsPercentage,
			node.Ip,
		)
	}
}
