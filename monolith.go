package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/blockc0de/monolith/internal/codes"
	"github.com/blockc0de/monolith/internal/config"
	"github.com/blockc0de/monolith/internal/handler"
	"github.com/blockc0de/monolith/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/monolith-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Authorization, Secret-Key")
	}))
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *codes.CodeError:
			return e.Code, e.Data()
		default:
			ex := codes.NewCodeError(http.StatusInternalServerError, e.Error()).Data()
			return http.StatusInternalServerError, ex
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
