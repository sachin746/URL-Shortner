package main

import (
	"context"
	"fmt"

	"URL-Shortner/api"
	"URL-Shortner/flags"
	"URL-Shortner/log"
	"URL-Shortner/utils/configs"
	"URL-Shortner/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	log.InitLogger(ctx)
	configs.InitConfigs(ctx)
	database.InitDatabase(ctx)
	startRouter()
}

func startRouter() {
	router, err := api.GetRouter()
	if err != nil {
		log.Sugar.Fatalf("Failed to get router: %v", err)
	}
	router.Use(gin.Recovery())
	log.Sugar.Info("Router initialized")

	err = router.Run(fmt.Sprintf(":%d", flags.Port()))
	if err != nil {
		log.Sugar.Fatalf("Failed to start server: %v", err)
	}
}
