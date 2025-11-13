package main

import (
	"context"
	"fmt"
	"habit-tracker/api"
	"habit-tracker/flags"
	"habit-tracker/log"
	"habit-tracker/utils"
	"habit-tracker/utils/configs"
	"habit-tracker/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	log.InitLogger(ctx)
	utils.InitValidator()
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
