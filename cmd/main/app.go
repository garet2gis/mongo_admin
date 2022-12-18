package main

import (
	"context"
	"log"
	"mongo_admin/internal/app"
	"mongo_admin/internal/config"
	"mongo_admin/pkg/client/mongodb"
)

func main() {
	cfg := config.GetConfig()
	ctx := context.Background()
	client, err := mongodb.NewClient(ctx, cfg.MongoDB.Host, cfg.MongoDB.Port)
	if err != nil {
		log.Fatal(err)
	}
	newApp := app.NewApp(ctx, client)
	newApp.Start()
}
