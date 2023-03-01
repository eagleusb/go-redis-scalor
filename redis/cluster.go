package rediscluster

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func rebalanceRedisShard(master string, slots int, fromId string, toId string) {
	cmd := []string{"--cluster", "reshard", master, "--cluster-from", fromId, "--cluster-to", toId, "--cluster-slots", "100", "--cluster-yes"}
	execRedisCli(cmd)
}

func (c *RedisClusterConf) WantedRedisConfArg(lookup string) bool {
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

func (c *RedisClusterConf) SetRedisConfArg(arg []string) {
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

func (c *RedisClusterConf) GetRedisConfArg(arg string) (_value string) {
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

func (r *RedisNodes) SetRedisClusterNodes(client *redis.Client, ctx context.Context) {
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

func (r *RedisNodes) GetRedisClusterNodes(client *redis.Client) {
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
