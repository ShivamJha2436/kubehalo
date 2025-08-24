package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ShivamJha2436/kubehalo/internal/kube"
	"github.com/ShivamJha2436/kubehalo/controllers/scalepolicy"
)

func main() {
	// Create clients
	_, dyn, _, err := kube.NewClients()
	if err != nil {
		log.Fatalf("failed to build clients: %v", err)
	}

	// Initialize lister
	lister := scalepolicy.NewLister(dyn)

	// Start Gin server
	r := gin.Default()

	// GET /scalepolicies
	r.GET("/scalepolicies", func(c *gin.Context) {
		ctx := context.Background()
		items, err := lister.ListScalePolicies(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert unstructured to plain JSON
		var results []map[string]interface{}
		for _, item := range items {
			results = append(results, item.Object)
		}
		c.JSON(http.StatusOK, results)
	})

	log.Println("[api] Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
