package logic

import (
	"context"
	"net/http"

	"github.com/blockc0de/monolith/internal/utils"

	"github.com/blockc0de/engine/compress"
	"github.com/blockc0de/monolith/internal/codes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"

	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DecompressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDecompressLogic(ctx context.Context, svcCtx *svc.ServiceContext) DecompressLogic {
	return DecompressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DecompressLogic) Decompress(req types.DecompressRequest) (resp *types.DecompressResponse, err error) {
	data, err := compress.GraphCompression{}.DecompressGraphData(req.Data)
	if err != nil {
		return nil, codes.NewCodeError(http.StatusBadRequest, "invalid graph")
	}

	projectId := gjson.Get(string(data), "project_id").String()
	address := common.HexToAddress(l.ctx.Value("address").(string))
	hash := utils.GetUniqueGraphHash(address, projectId)

	return &types.DecompressResponse{Decompressed: string(data), Hash: hash}, nil
}
