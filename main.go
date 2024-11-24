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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(cors.New(corsConfig))
	
	server.GET("/blocks", api.GetBlocks)
	server.POST("/blocks", api.AddBlock)
	server.POST("/upload-file-block", api.AddBlockFile)
	server.POST("/sync", api.Synchronize)
	server.POST("/validate", api.ValidateBlock)
	

	listenPort := ":" + os.Args[1] + "0"


	server.Run(listenPort) 

}


