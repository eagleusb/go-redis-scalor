package cmd

import (
	"fmt"

	"github.com/eagleusb/go-redis-scalor/utilities"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current redis cluster status with slots information",
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
		// print current redis shards list
		if clusterNodesList, err := client.ClusterNodes(ctx).Result(); err == nil {
			fmt.Printf(
				"Current Cluster Nodes List:\n%s\n",
				clusterNodesList,
			)
		}
	},
}
