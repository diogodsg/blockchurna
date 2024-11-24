package api

import (
	"blockchurna/blockchain"
	"blockchurna/p2p"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)



func GetBlocks(c *gin.Context) {
	data := blockchain.BC.AggregateVotes()
	print(data)
	c.JSON(200, gin.H{
		"blocks": blockchain.BC.Blocks,
		"aggregated": blockchain.BC.AggregateVotes(),
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

	createdBlock, err := blockchain.BC.CreateBlock(createBlockDto)
	if err != nil {
		c.JSON(200, gin.H{
			"error":"error",
		})
	}
	p2p.BroadcastBlock(p2p.Node.Host, createdBlock)
	c.JSON(200, gin.H{
		"block":createdBlock,
	})
}

func AddBlockFile(c *gin.Context) {
		// Receive the file from the request
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		defer file.Close()
	
		// Read the file content into a byte slice
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
	
		// Unmarshal the JSON into a Block struct
		var block blockchain.Payload
		err = json.Unmarshal(fileContent, &block)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao carregar bloco"})
			return
		}		
		createdBlock, err := blockchain.BC.CreateBlock(block)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p2p.BroadcastBlock(p2p.Node.Host, createdBlock)

		// Respond with a success message
		c.JSON(http.StatusOK, gin.H{"message": "Bloco adicionado com sucesso"})
	
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

