package hashutil

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/agungangga/corsair-blockchain/blockdata"
)

func GetBlockHash(block blockdata.Block) string {
	block.Hash = ""
	block.PrevHash = ""

	jBlock, _ := json.Marshal(block)

	b := md5.New().Sum(jBlock)
	return hex.EncodeToString(b)
}
