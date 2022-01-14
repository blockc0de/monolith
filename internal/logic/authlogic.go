package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/blockc0de/monolith/internal/svc"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	EthereumSignKey  = "I agree to connect my wallet to the GraphLinq Interface."
	EthereumSignHash = crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(EthereumSignKey), EthereumSignKey)))
)

type AuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) AuthLogic {
	return AuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthLogic) Auth(req types.AuthRequest) (resp *types.AuthResponse, err error) {
	sig, err := hexutil.Decode(req.Signature)
	if err != nil {
		return &types.AuthResponse{Auth: false}, nil
	}
	if sig[64] != 27 && sig[64] != 28 {
		return &types.AuthResponse{Auth: false}, nil
	}
	sig[64] -= 27

	sigPublicKey, err := crypto.SigToPub(EthereumSignHash, sig)
	if err != nil {
		return &types.AuthResponse{Auth: false}, nil
	}

	address := crypto.PubkeyToAddress(*sigPublicKey)
	if address != common.HexToAddress(req.Address) {
		return &types.AuthResponse{Auth: false}, nil
	}

	now := time.Now().Unix()
	config := l.svcCtx.Config.Auth
	accessToken, err := l.getJwtToken(config.AccessSecret, now, config.AccessExpire, address.String())
	if err != nil {
		return &types.AuthResponse{Auth: false}, nil
	}

	return &types.AuthResponse{Auth: true, AccessToken: &accessToken}, nil
}

func (l *AuthLogic) getJwtToken(secretKey string, iat, seconds int64, address string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["address"] = address
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
