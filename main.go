package main

import (
	"blockchurna/api"
	"blockchurna/p2p"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	p2p.StartBlockchain()

	server := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(cors.New(corsConfig))
	
	server.GET("/core/v1/blocks", api.GetBlocks)
	server.POST("/core/v1/blocks", api.AddBlock)
	server.POST("/core/v1/verify-vote", api.VerifyVote)
	server.POST("/core/v1/upload-file-block", api.AddBlockFile)
	server.POST("/core/v1/sync", api.Synchronize)
	server.POST("/core/v1/validate", api.ValidateBlock)
	

	listenPort := ":" + "40010"


	server.Run(listenPort) 
}


