package jirametrics

import (
	"log"
)

func UpdateMetrics(config *Config, username string, token string, dryrun bool) error {

	// process custom metrics
	// only report total number of matches
	for index, metric := range config.Metrics {

		if dryrun {
			continue
		}

		result, err := SimpleQuery(config, username, token, dryrun, metric.Query)

		if err != nil {
			return err
		}

		log.Printf("[custom] %s %d\n", metric.Query, result.Total)

		switch metric.MetricType {
		case "counter":
			config.Metrics[index].Counter.Add(float64(result.Total))
		case "gauge":
			config.Metrics[index].Gauge.Set(float64(result.Total))
		case "histogram":
			config.Metrics[index].Histogram.Observe(float64(result.Total))
		}
	}

	if dryrun {
		return nil
	}

	// reset metrics as required before iterating over queues
	config.NewIncidents15mGauge.Reset()
	config.NewTickets15mGauge.Reset()

	err := processQueues(config, username, token, dryrun)
	if err != nil {
		return err
	}

	err = processBoards(config, username, token, dryrun)
	if err != nil {
		return err
	}
	return nil
}
