package dsncli

import (
	"context"
	"github.com/Bedrock-Technology/Dsn/app"
	"github.com/Bedrock-Technology/Dsn/app/api"
	"github.com/Bedrock-Technology/Dsn/app/dsn"
	"github.com/Bedrock-Technology/Dsn/app/node"
	_ "github.com/Bedrock-Technology/Dsn/docs"
	"github.com/Bedrock-Technology/Dsn/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// @title           Web Server API
// @version         1.0
// @description     This is webapp server api.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  hanson@bedrock.technology

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a dns process",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := dsn.GetConfig()
		e := gin.Default()
		setupRouter(e)
		setupSwagger(e, conf.DocAuth)
		gin.SetMode(gin.ReleaseMode)

		srv := &http.Server{
			Addr:    conf.Rpc,
			Handler: e,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Errorf("listen and serve", "err", err)
				panic(err)
			}
		}()

		shutdownChan := make(chan struct{})
		stopChan := make(chan struct{})
		shutdown := node.StopFunc(func(ctx context.Context) error {
			close(stopChan)
			log.Info("toml monitor shutdown")
			return nil
		})
		app.GetConfigStore().RunWithCron(conf.TomlDir, conf.GitUpdateShell, conf.CheckSpec, conf.ReloadFlag, stopChan)

		configFile, _ := cmd.Root().PersistentFlags().GetString("config")
		dsn.RunWithWatch(configFile, stopChan)

		// Monitor for shutdown.
		finishCh := node.MonitorShutdown(shutdownChan, node.ShutdownHandler{Component: "toml-monitor", StopFunc: shutdown})
		<-finishCh

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorf("dns", "Server Shutdown:", err)
		}
		log.Info("Server exiting")
		return nil
	},
}

func setupRouter(e *gin.Engine) {
	root := e.Group("/")
	{
		r := root.Group("dsn")
		_api := api.DsnApi{}
		r.GET("/exec/:func_name/:params", _api.ExecCmd)
		r.GET("/func_list", _api.GetFuncList)
		r.POST("/load", _api.LoadSqlFile)
		r.POST("/execsql", _api.ExecSql)
	}
}

func setupSwagger(e *gin.Engine, docAuth map[string]string) {
	auth := gin.Accounts(docAuth)
	e.GET("/docs/*any", gin.BasicAuth(auth), gs.WrapHandler(sf.Handler))
	e.StaticFile("/swagger/doc.json", "./docs/swagger.json")
}
