package dsncli

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/http"
	"github.com/spf13/cobra"
)

var SqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "load sql config file to the dsn app",
	RunE: func(cmd *cobra.Command, args []string) error {
		rpc, _ := cmd.Flags().GetString("rpc")
		tomlFile, _ := cmd.Flags().GetString("toml-file")
		fmt.Println("rpc is", rpc)
		rep, err := http.LoadSqlConfig(rpc, tomlFile)
		fmt.Println("load sql config file to the dsn app", rep, "err is ", err)
		return err
	},
}

func init() {
	SqlCmd.Flags().String("rpc", "http://127.0.0.1:1234/dsn/load", "local rpc server address")
	SqlCmd.Flags().String("toml-file", "sql.toml", "add the sql config file to the dsn app")
}
