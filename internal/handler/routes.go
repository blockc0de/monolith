// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/blockc0de/monolith/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/wallets/auth",
				Handler: AuthHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/graphs/template/:id",
				Handler: GetGraphsTemplateHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/graphs/logs",
				Handler: LogsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/graphs/deploy",
				Handler: DeployHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/graphs/compress",
				Handler: CompressHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/graphs/decompress",
				Handler: DecompressHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
