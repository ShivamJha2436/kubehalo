package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ShivamJha2436/kubehalo/controllers/scalepolicy"
	"github.com/ShivamJha2436/kubehalo/internal/config"
	"github.com/ShivamJha2436/kubehalo/internal/kube"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadAPIConfig()

	clients, err := kube.NewClients()
	if err != nil {
		log.Fatalf("failed to build clients: %v", err)
	}

	lister := scalepolicy.NewLister(clients.Dynamic)

	r := gin.Default()

	r.GET("/scalepolicies", func(c *gin.Context) {
		ctx := context.Background()
		items, err := lister.ListScalePolicies(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []map[string]interface{}
		for _, item := range items {
			results = append(results, item.Object)
		}
		c.JSON(http.StatusOK, results)
	})

	log.Printf("[api] Starting server on %s", cfg.Address)
	if err := r.Run(cfg.Address); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
