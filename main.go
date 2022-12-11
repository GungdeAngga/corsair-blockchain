package main

import (
	"strconv"

	"github.com/agungangga/corsair-blockchain/blockchainutil"
	"github.com/agungangga/corsair-blockchain/blockdata"
	"github.com/agungangga/corsair-blockchain/blockmanager"
	"github.com/agungangga/corsair-blockchain/transferdata"
	"github.com/agungangga/corsair-blockchain/userdata"
	"github.com/labstack/echo/v4"
)

func findDefaultBlockchainWithConsensus() blockdata.Blockchain {
	blockchain1 := blockmanager.GetBlockchain(1)
	blockchain2 := blockmanager.GetBlockchain(2)
	blockchain3 := blockmanager.GetBlockchain(2)

	blockchain := make(blockdata.Blockchain, 0)
	longestBlockchain := blockchainutil.FindLongestBlockchainLength(
		blockchain1,
		blockchain2,
		blockchain3,
	)

	for i := 0; i < longestBlockchain; i++ {
		blocksToSearch := make(blockdata.Blockchain, 0)

		if i < len(blockchain1) {
			blocksToSearch = append(blocksToSearch, blockchain1[i])
		}

		if i < len(blockchain2) {
			blocksToSearch = append(blocksToSearch, blockchain2[i])
		}

		if i < len(blockchain3) {
			blocksToSearch = append(blocksToSearch, blockchain3[i])
		}

		highestConsensusBlock := blockchainutil.FindConsensus(
			blocksToSearch...,
		)
		blockchain = append(blockchain, highestConsensusBlock)
	}
	return blockchain
}

func main() {
	ec := echo.New()

	ec.GET("/blockchain", func(c echo.Context) error {
		return c.JSON(200, findDefaultBlockchainWithConsensus())
	})

	ec.POST("/blockchain/search-by-hash", func(c echo.Context) error {
		var hashSearch blockdata.SearchByHash
		if err := c.Bind(&hashSearch); err != nil {
			return c.String(500, err.Error())
		}

		blockchains := findDefaultBlockchainWithConsensus()
		var chosenBlockchain *blockdata.Block

		for _, v := range blockchains {
			if v.Hash == hashSearch.Hash {
				chosenBlockchain = &v
			}
		}

		if chosenBlockchain == nil {
			return c.String(404, "block tidak ditemukan")
		}

		return c.JSON(200, chosenBlockchain)
	})

	ec.GET("/transaction/:id", func(c echo.Context) error {
		idStr := c.Param("id")

		id64, _ := strconv.ParseInt(idStr, 10, 32)

		id := int(id64)

		blockchains := findDefaultBlockchainWithConsensus()
		filteredBlocks := make(blockdata.Blockchain, 0)

		for _, v := range blockchains {
			if v.SenderID == id || v.ReceiverID == id {
				filteredBlocks = append(filteredBlocks, v)
			}
		}

		return c.JSON(200, filteredBlocks)
	})

	ec.GET("/account/:id", func(c echo.Context) error {
		idStr := c.Param("id")

		id64, _ := strconv.ParseInt(idStr, 10, 32)

		id := int(id64)

		user := userdata.User{
			ID: id,
		}

		blockchains := findDefaultBlockchainWithConsensus()
		for _, v := range blockchains {
			if v.ReceiverID == id {
				user.Balance += v.Amount
				continue
			}
			if v.SenderID == id {
				user.Balance -= v.Amount
				continue
			}
		}

		return c.JSON(200, user)
	})

	ec.POST("/transfer", func(c echo.Context) error {
		var trfData transferdata.TransferData
		if err := c.Bind(&trfData); err != nil {
			return c.String(500, err.Error())
		}

		var senderBalance int64 = 0
		blockchains := findDefaultBlockchainWithConsensus()
		for _, v := range blockchains {
			if v.ReceiverID == trfData.From {
				senderBalance += v.Amount
				continue
			}
			if v.SenderID == trfData.From {
				senderBalance -= v.Amount
				continue
			}
		}

		if senderBalance < trfData.Amount {
			return c.String(400, "dana tidak mencukupi")
		}

		newBlock := blockmanager.NewBlock(trfData.From, trfData.To, trfData.Amount, "")
		hash := blockmanager.PresistBlock(1, newBlock)
		blockmanager.PresistBlock(2, newBlock)
		blockmanager.PresistBlock(3, newBlock)

		newBlock.PrevHash = hash

		return c.JSON(200, newBlock)
	})

	ec.Start("localhost:8080")
}
