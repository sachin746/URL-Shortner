package main

import (
	"context"
	"fmt"
	"time"

	"URL-Shortner/api"
	"URL-Shortner/flags"
	"URL-Shortner/log"
	"URL-Shortner/utils/cache"
	"URL-Shortner/utils/configs"
	"URL-Shortner/utils/cron"
	"URL-Shortner/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	log.InitLogger(ctx)
	configs.InitConfigs(ctx)
	cache.InitRedisCache(ctx)
	database.InitDatabase(ctx)
	cron.StartCleanupCron(ctx, 1*time.Hour)
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
