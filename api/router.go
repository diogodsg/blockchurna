package api

import (
	"blockchurna/blockchain"
	"blockchurna/p2p"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBlockDto struct {
	Payload string
}

func GetBlocks(c *gin.Context) {
	c.JSON(200, gin.H{
		"blocks":blockchain.BC.Blocks,
	})
}


func AddBlock(c *gin.Context) {
	var createBlockDto CreateBlockDto
	err := c.ShouldBindJSON(&createBlockDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse data"})
		return
	}

	createdBlock := blockchain.BC.CreateBlock(createBlockDto.Payload)
	p2p.BroadcastBlock(p2p.Node.Host, createdBlock)
	c.JSON(200, gin.H{
		"block":createdBlock,
	})
}

func Synchronize(c *gin.Context) {
	p2p.SynchronizeChain(p2p.Node.Host,blockchain.BC)
	c.JSON(200, gin.H{
		"block":"createdBlock",
	})
}