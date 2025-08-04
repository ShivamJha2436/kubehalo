package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"kubehalo/api/v1"
)

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func FetchMetric(policy v1.ScalePolicy) (float64, error) {
	url := fmt.Sprintf("http://prometheus:9090/api/v1/query?query=%s", policy.Spec.MetricQuery)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Data.Result) == 0 {
		return 0, fmt.Errorf("no result from Prometheus")
	}

	value := result.Data.Result[0].Value[1].(string)
	var val float64
	fmt.Sscanf(value, "%f", &val)
	return val, nil
}
