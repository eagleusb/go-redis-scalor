package main

import (
	"context"
	"fmt"
	"strings"
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

type RedisNodes struct {
	ClusterNodes map[string]RedisNode
}

type RedisClusterConf struct {
	Conf          map[string]string
	State         string
	Size          string
	SlotsAssigned string
	SlotsOk       string
	NodesKnown    string
}

func init() { fmt.Print("Hello World\n") }

func main() {

	// TODO: flags

	clientOpts := &redis.Options{
		Addr:        ":32000",
		ClientName:  "go-redis-scalor",
		DialTimeout: timeout,
	}
	client := redis.NewClient(clientOpts)

	if client != nil {
		fmt.Printf("Connected to: %s\n", client.Options().Addr)
		_, err := client.Ping(ctx).Result()
		checkErr(err)
	}

	// NOTE: not implemented by redis lib defer client.Quit(ctx)

	// print current redis shards list
	if clusterNodesList, err := client.ClusterNodes(ctx).Result(); err == nil {
		fmt.Printf(
			"Current Cluster Nodes List:\n%s\n",
			clusterNodesList,
		)
	}

	// gather redis shards with slots information
	clusterNodes := &RedisNodes{
		ClusterNodes: make(map[string]RedisNode),
	}
	clusterNodes.setRedisClusterNodes(client)
	clusterNodes.getRedisClusterNodes(client)


	// gather redis cluster configuration
	clusterConf := &RedisClusterConf{
		Conf: make(map[string]string),
	}
	clusterInfo, _ := client.ClusterInfo(ctx).Result()

	for _, v := range strings.Split(clusterInfo, "\n") {
		_args := strings.Split(v, ":")
		if len(_args) >= 2 {
			clusterConf.Conf[_args[0]] = _args[1]

			if clusterConf.wantedRedisConfArg(_args[0]) {
				clusterConf.setRedisConfArg(_args)
			}
		}
	}

	/*
		TODO: live resharding
		implement live resharding through API
		see https://redis.io/commands/cluster-setslot/
	*/
	// rebalanceRedisShard("0:32000", clusterNodes["53ed71dd58d1f04de0a08a22e9ea686a3ae6599b"].SlotsCount, "53ed71dd58d1f04de0a08a22e9ea686a3ae6599b", "3578e289c1cef4099f345b284fe362f990f6a76e")

}
