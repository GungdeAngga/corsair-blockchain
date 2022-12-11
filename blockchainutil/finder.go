package blockchainutil

import (
	"github.com/agungangga/corsair-blockchain/blockdata"
	"github.com/agungangga/corsair-blockchain/hashutil"
)

func FindLongestBlockchainLength(blockchains ...blockdata.Blockchain) int {
	var max int = 0
	for _, v := range blockchains {
		if len(v) > max {
			max = len(v)
		}
	}
	return max
}

func FindConsensus(blocks ...blockdata.Block) blockdata.Block {
	consensusCount := make(map[string]int)
	var highestBlock blockdata.Block
	highest := 0

	for _, v := range blocks {
		hash := hashutil.GetBlockHash(v)
		consensusCount[hash]++
		if consensusCount[hash] > highest {
			highestBlock = v
		}
	}

	return highestBlock
}
