package logic

import (
	"context"
	"net/http"

	"github.com/blockc0de/monolith/internal/codes"

	"github.com/blockc0de/monolith/internal/storage"

	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogsLogic {
	return LogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogsLogic) Logs(req types.LogsRequest) (resp *types.LogsResponse, err error) {
	graphs := storage.GraphsManager{RedisClient: l.svcCtx.RedisClient}
	logs, err := graphs.GetLogs(req.Hash)
	if err != nil {
		return nil, codes.NewCodeError(http.StatusInternalServerError, "internal server error")
	}

	return &types.LogsResponse{Logs: logs}, nil
}
