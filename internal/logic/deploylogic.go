package logic

import (
	"context"

	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeployLogic {
	return DeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeployLogic) Deploy(req types.DeployRequest) (resp *types.DeployResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
