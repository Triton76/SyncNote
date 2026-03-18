// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"SyncNote/pkg/httpx"
	"SyncNote/syncnote/api/internal/config"
	"SyncNote/syncnote/api/internal/handler"
	"SyncNote/syncnote/api/internal/svc"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/syncnote-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(httpx.CorsMiddleware())

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
