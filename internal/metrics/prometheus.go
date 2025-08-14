package metrics

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
)

type PrometheusClient struct {
	client v1.API
}

func NewPrometheusClient(address string) (*PrometheusClient, error) {
	cfg := api.Config{Address: address}
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &PrometheusClient{client: v1.NewAPI(client)}, nil
}

// QueryMetric executes a PromQL query and returns the first value
func (p *PrometheusClient) QueryMetric(query string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := p.client.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		fmt.Println("Prometheus warnings:", warnings)
	}

	// Extract the value (assumes scalar result)
	vector, ok := result.(v1.Vector)
	if ok && len(vector) > 0 {
		return float64(vector[0].Value), nil
	}

	return 0, fmt.Errorf("no data returned for query: %s", query)
}
