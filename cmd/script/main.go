package main

import (
	"context"
	"flag"
	"log"

	"omiai-server/cmd/script/command"
	"omiai-server/internal/conf"

	"github.com/google/wire"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                "script",
		DisableFlagParsing: true,
	}

	// confPath is the conf flag.
	confPath   string
	confCenter string
)

func init() {
	flag.StringVar(&confPath, "conf", "configs", "conf path, eg: -conf conf.yaml")
	flag.StringVar(&confCenter, "c", "", "config center, eg: -c nacos://nacos-prod.wxbjq.net.cn?namespace_id=b265887d-5811-478a-8512-2a9462ce0431&data_id=aippt-server.yaml&group=aippt-server-dev&log_level=error")
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(InitCmd), "*"),
)

type InitCmd struct {
	Command *command.Script
}

func main() {
	flag.Parse()
	var err error
	_, err = conf.Init(confPath, confCenter)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// initApp
	app, cleanup, err := initApp(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Add commands
	rootCmd.AddCommand(app.Command.InsertClass())
	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("execute core service failed, %s", err.Error())
	}
}
