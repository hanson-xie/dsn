package main

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/cmd/dsncli"
	"os"

	"github.com/spf13/cobra"

	"github.com/Bedrock-Technology/Dsn/app/dsn"
	"github.com/Bedrock-Technology/Dsn/build"
	"github.com/Bedrock-Technology/Dsn/log"
)

var (
	configFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "dsn-cli",
		Short:   "dsn client app control center",
		Version: build.UserVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := dsn.LoadConfig(configFile)
			if err != nil {
				fmt.Printf("failed load config: %v\n", err)
				os.Exit(1)
			}
			conf := dsn.GetConfig()
			log.ConfigLog(&conf.Log)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if r := recover(); r != nil {
				// Generate report in BEDROCK_REWARD_PATH and re-raise panic
				panic(r)
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "app.yaml", "the config file of btc redeem app")

	rootCmd.AddCommand(dsncli.RunCmd)
	rootCmd.AddCommand(dsncli.SqlCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("app run error: %v\n", err)
		os.Exit(1)
	}
}
