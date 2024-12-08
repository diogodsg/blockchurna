package main

import (
	"blockchurna/api"
	"blockchurna/p2p"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	p2p.StartBlockchain()

	server := gin.Default()

	

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000", "http://blockchurna.tech"}, // Allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	
	server.GET("/blocks", api.GetBlocks)
	server.POST("/blocks", api.AddBlock)
	server.POST("/verify-vote", api.VerifyVote)
	server.POST("/upload-file-block", api.AddBlockFile)
	server.POST("/sync", api.Synchronize)
	server.POST("/validate", api.ValidateBlock)
	

	listenPort := ":" + os.Args[1] + "0"


	server.Run(listenPort) 

}


