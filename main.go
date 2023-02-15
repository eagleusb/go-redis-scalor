package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx context.Context = context.Background()
var timeout time.Duration = 30 * time.Second

type RedisNode struct {
	Slots           []string
	SlotsCount      int
	SlotsPercentage int
	Id              string
	Ip              string
	Name            string
}

func init() { fmt.Print("Hello World\n") }

func main() {

	var clientOpts = &redis.Options{
		Addr:        ":6379",
		ClientName:  "go-redis-scalor",
		DialTimeout: timeout,
	}
	client := redis.NewClient(clientOpts)
	// TODO defer client.Quit(ctx)

	fmt.Printf("Connected to: %s\n", client.Options().Addr)

	_, err := client.Ping(ctx).Result()
	checkErr(err)

	clusterInfo, _ := client.ClusterInfo(ctx).Result()
	fmt.Printf(
		"Current Cluster State:\n%s\n",
		clusterInfo,
	)

	clusterNodesList, _ := client.ClusterNodes(ctx).Result()
	fmt.Printf(
		"Current Cluster Nodes List:\n%s\n",
		clusterNodesList,
	)

	clusterSlots, _ := client.ClusterSlots(ctx).Result()
	clusterNodes := make(map[string]RedisNode)

	for i, slot := range clusterSlots {
		_nodeId := clusterSlots[i].Nodes[0].ID
		_nodeIp := clusterSlots[i].Nodes[0].Addr

		clusterNode := &RedisNode{
			Id: _nodeId,
			Ip: _nodeIp,
		}

		if _, ok := clusterNodes[_nodeId]; ok {
			clusterNode.Slots = append(clusterNodes[_nodeId].Slots, fmt.Sprint(slot.Start, "-", slot.End))
			clusterNode.SlotsCount = clusterNodes[_nodeId].SlotsCount + slotsCount(slot.Start, slot.End)
			clusterNode.SlotsPercentage = clusterNodes[_nodeId].SlotsPercentage + slotsPercentage(slotsCount(slot.Start, slot.End))
		} else {
			clusterNode.Slots = make([]string, 1)
			clusterNode.Slots[0] = fmt.Sprint(slot.Start, "-", slot.End)
			clusterNode.SlotsCount = slotsCount(slot.Start, slot.End)
			clusterNode.SlotsPercentage = slotsPercentage(clusterNode.SlotsCount)
		}

		clusterNodes[_nodeId] = *clusterNode
	}

	for _, node := range clusterNodes {
		fmt.Printf(
			"ID: %s slots: %s slotsEntries: %d slotsCount: %d slotsPercentage: %d IP: %s\n",
			node.Id,
			node.Slots,
			len(node.Slots),
			node.SlotsCount,
			node.SlotsPercentage,
			node.Ip,
		)
	}

	/*
		TODO: implement live resharding through API
		see https://redis.io/commands/cluster-setslot/
	*/
	// v, err := client.Do(ctx, "cluster", "setslot"...).Text()
	// fmt.Printf("%q %s", v, err)

	execRedisCli("--version")
}
