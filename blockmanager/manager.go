package blockmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/agungangga/corsair-blockchain/blockdata"
	"github.com/agungangga/corsair-blockchain/hashutil"
)

func NewBlock(SenderID int, ReceiverID int, amount int64, prevHash string) blockdata.Block {
	timestamp := time.Now().UnixMilli()
	b := blockdata.Block{
		SenderID:   SenderID,
		ReceiverID: ReceiverID,
		Amount:     amount,
		PrevHash:   prevHash,
		Timestamp:  timestamp,
	}

	hashValue := hashutil.GetBlockHash(b)
	b.Hash = hashValue
	return b
}

func PresistBlock(blockNode int, block blockdata.Block) string {
	blockchain := GetBlockchain(blockNode)
	var prevHash string

	if len(blockchain) > 0 {
		prevHash = blockchain[len(blockchain)-1].Hash
	}

	block.PrevHash = prevHash
	blockchain = append(blockchain, block)

	jsonified, _ := json.Marshal(blockchain)

	err := os.WriteFile(fmt.Sprintf("./blockchain-%d.json", blockNode), jsonified, 0644)
	if err != nil {
		log.Println(err)
	}
	return block.PrevHash
}

func GetBlockchain(blockNode int) blockdata.Blockchain {
	blockchain := make(blockdata.Blockchain, 0)
	f, err := os.ReadFile(fmt.Sprintf("./blockchain-%d.json", blockNode))
	if err != nil {
		log.Println(err)
		return blockchain
	}

	err = json.Unmarshal(f, &blockchain)
	if err != nil {
		log.Println(err)
		return blockchain
	}
	return blockchain
}
