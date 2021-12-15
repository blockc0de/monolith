package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func GetUniqueGraphHash(address common.Address, projectId string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", address.String(), projectId)))
	return hex.EncodeToString(sum[:])
}
