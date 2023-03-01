package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var ctx context.Context = context.Background()
var timeout time.Duration = 30 * time.Second

var rootCmd = &cobra.Command{
  Use:   "redis-scalor",
  Short: "Redis scaler: migrate slots from/to shards, get the status of the cluster",
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println("Error while executing the cli", err)
    os.Exit(1)
  }
}
