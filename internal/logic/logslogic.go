package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
