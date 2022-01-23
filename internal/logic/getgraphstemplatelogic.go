package logic

import (
	"context"

	"github.com/blockc0de/monolith/internal/storage"
	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGraphsTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGraphsTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGraphsTemplateLogic {
	return GetGraphsTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGraphsTemplateLogic) GetGraphsTemplate(req types.GetGraphsTemplateRequest) (resp *types.GetGraphsTemplateResponse, err error) {
	resp = &types.GetGraphsTemplateResponse{}

	templates := storage.TemplateManager{RedisClient: l.svcCtx.RedisClient}
	template, err := templates.Get(req.Id)
	if err != nil {
		resp.Success = false
		return resp, nil
	}

	resp.Template = &types.GraphTemplate{
		IdGraphsTemplates: req.Id,
		Title:             template.Title,
		Key:               template.Key,
		Bytes:             template.Bytes,
		Description:       template.Description,
		CustomImg:         template.CustomImg,
	}

	return resp, nil
}
