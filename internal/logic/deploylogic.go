package logic

import (
	"context"
	"net/http"

	"github.com/blockc0de/engine/interop"

	"github.com/blockc0de/engine/compress"
	"github.com/blockc0de/monolith/internal/codes"
	"github.com/blockc0de/monolith/internal/storage"
	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/blockc0de/monolith/internal/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
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
	data, err := compress.GraphCompression{}.DecompressGraphData(req.Bytes)
	if err != nil {
		return nil, codes.NewCodeError(http.StatusBadRequest, "invalid graph")
	}

	graph, err := interop.LoadGraph(data)
	if err != nil {
		return nil, codes.NewCodeError(http.StatusBadRequest, "invalid graph")
	}

	projectId := gjson.Get(string(data), "project_id").String()
	address := common.HexToAddress(l.ctx.Value("address").(string))
	hash := utils.GetUniqueGraphHash(address, projectId)

	graphs := storage.GraphsManager{RedisClient: l.svcCtx.RedisClient}
	if err = graphs.Save(hash, address.String(), req.Bytes); err != nil {
		return nil, codes.NewCodeError(http.StatusInternalServerError, "internal server error")
	}
	graphs.ClearLogs(hash)

	wallets := storage.WalletsManager{RedisClient: l.svcCtx.RedisClient}
	if err = wallets.AddGraph(address.String(), hash); err != nil {
		return nil, codes.NewCodeError(http.StatusInternalServerError, "internal server error")
	}

	graph.Hash = hash
	l.svcCtx.GraphContainer.AddNewGraph(address, hash, graph)
	logx.Infof("Graph deployed, wallet: %s, hash: %s", address.String(), hash)

	return &types.DeployResponse{Hash: hash}, nil
}
