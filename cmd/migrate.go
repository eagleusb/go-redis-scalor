package cmd

import (
	"fmt"
	"strings"

	rediscluster "github.com/eagleusb/go-redis-scalor/redis"
	"github.com/eagleusb/go-redis-scalor/utilities"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate slots from one redis shard to another",
	Run: func(cmd *cobra.Command, args []string) {
		clientOpts := &redis.Options{
			Addr:        ":32000",
			ClientName:  "go-redis-scalor",
			DialTimeout: timeout,
		}
		client := redis.NewClient(clientOpts)

		if client != nil {
			fmt.Printf("Connected to: %s\n", client.Options().Addr)
			_, err := client.Ping(ctx).Result()
			utilities.CheckErr(err)
		}

		// NOTE: not implemented by redis lib defer client.Quit(ctx)

		// gather redis shards with slots information
		clusterNodes := &rediscluster.RedisNodes{
			ClusterNodes: make(map[string]rediscluster.RedisNode),
		}
		clusterNodes.SetRedisClusterNodes(client, ctx)
		clusterNodes.GetRedisClusterNodes(client)

		// gather redis cluster configuration
		clusterConf := &rediscluster.RedisClusterConf{
			Conf: make(map[string]string),
		}
		clusterInfo, _ := client.ClusterInfo(ctx).Result()

		for _, v := range strings.Split(clusterInfo, "\n") {
			_args := strings.Split(v, ":")
			if len(_args) >= 2 {
				clusterConf.Conf[_args[0]] = _args[1]

				if clusterConf.WantedRedisConfArg(_args[0]) {
					clusterConf.SetRedisConfArg(_args)
				}
			}
		}

		/*
			TODO: live resharding
			implement live resharding through API
			see https://redis.io/commands/cluster-setslot/
		*/
		// rebalanceRedisShard("0:32000", clusterNodes["53ed71dd58d1f04de0a08a22e9ea686a3ae6599b"].SlotsCount, "53ed71dd58d1f04de0a08a22e9ea686a3ae6599b", "3578e289c1cef4099f345b284fe362f990f6a76e")

	},
}
