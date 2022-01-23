package logic

import (
	"context"
	"net/http"

	"github.com/blockc0de/engine/compress"
	"github.com/blockc0de/engine/interop"
	"github.com/blockc0de/monolith/internal/codes"
	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/blockc0de/monolith/internal/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

type CompressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompressLogic(ctx context.Context, svcCtx *svc.ServiceContext) CompressLogic {
	return CompressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompressLogic) Compress(req types.CompressRequest) (resp *types.CompressResponse, err error) {
	compressed, err := compress.GraphCompression{}.CompressGraphData([]byte(req.Data))
	if err != nil {
		return nil, codes.NewCodeError(http.StatusBadRequest, "invalid graph")
	}

	_, err = interop.LoadGraph([]byte(req.Data))
	if err != nil {
		return nil, codes.NewCodeError(http.StatusBadRequest, "invalid graph")
	}

	projectId := gjson.Get(req.Data, "project_id").String()
	address := common.HexToAddress(l.ctx.Value("address").(string))
	hash := utils.GetUniqueGraphHash(address, projectId)

	return &types.CompressResponse{Compressed: compressed, Hash: hash}, nil
}
