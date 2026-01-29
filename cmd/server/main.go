package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/iWuxc/go-wit/queue/server"

	"omiai-server/internal/conf"
	"omiai-server/pkg/validate"

	"github.com/gin-gonic/gin/binding"
	"github.com/iWuxc/go-wit/app"
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/transport"
	"github.com/iWuxc/go-wit/validator"
	"github.com/libi/dcron"
)

var (
	// Version auto generate. DO NOT EDIT.
	// go build -ldflags "-X main.Version=v1.0.0"
	Version string
	// confPath is the conf flag.
	confPath   string
	confCenter string
)

func init() {
	flag.StringVar(&confPath, "conf", "configs", "conf path, eg: -conf conf.yaml")
	flag.StringVar(&confCenter, "c", "", "config center, eg: -c nacos://nacos.wxbjq.net.cn?namespace_id=b265887d-5811-478a-8512-2a9462ce0431&data_id=aippt-server.yaml&group=aippt-server-dev&log_level=error")
}

func newApp(ctx context.Context, hs []transport.ServerInterface, cron *dcron.Dcron, server *server.Server) (*app.App, func(), error) {
	log.Printf("App Version: %s", Version)
	// 自定义表单校验
	binding.Validator = validator.NewValidator(validate.NewValidator())

	go func() {
		if conf.GetConfig().Cron {
			cron.Start()
		}
	}()
	a, err := app.NewApp(app.Context(ctx), app.Server(hs...), app.Version(Version), app.Name("aicoloset-server"))
	return a, func() {
		func() {
			ctx := cron.Stop()
			select {
			case <-ctx.Done():
			case <-time.After(60 * time.Second):
				log.Info("context was not done immediately")
			}
		}()
		server.Shutdown()
		server.Stop()
	}, err
}

func main() {
	flag.Parse()
	f, err := conf.Init(confPath, confCenter)
	if err != nil {
		panic(err)
	}
	defer f()

	ctx := context.Background()

	// initApp
	a, cleanup, err := initApp(ctx)

	if err != nil {
		panic(err)
	}
	defer cleanup()

	if f, err := os.OpenFile(conf.GetConfig().Log.Path+"crash.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm); err == nil {
		defer f.Close()
		if e := crashLog(f); e != nil {
			panic(e)
		}
	}

	if e := a.Run(); e != nil {
		panic(e)
	}
}
