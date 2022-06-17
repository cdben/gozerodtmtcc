package main

import (
	"flag"
	"fmt"

	_ "github.com/dtm-labs/dtmdriver-gozero"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"tcc/api/internal/config"
	"tcc/api/internal/handler"
	"tcc/api/internal/svc"
)

var configFile = flag.String("f", "api/etc/trans.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
