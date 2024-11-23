package api

import (
	"blockchurna/blockchain"
	"blockchurna/p2p"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)



func GetBlocks(c *gin.Context) {
	c.JSON(200, gin.H{
		"blocks":blockchain.BC.Blocks,
	})
}


func AddBlock(c *gin.Context) {
	fmt.Println("afwafmkafmk")
	var createBlockDto blockchain.Payload
	err := c.ShouldBindJSON(&createBlockDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse data"})
		return
	}

	createdBlock := blockchain.BC.CreateBlock(createBlockDto)
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

func ValidateBlock(c *gin.Context) {
	var payload blockchain.Payload
	err := c.ShouldBindJSON(&payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse data"})
		return
	}

	 blockchain.ValidatePayload(payload)


}

