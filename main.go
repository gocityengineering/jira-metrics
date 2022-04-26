package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gocityengineering/jira-metrics/encodingutils"
	"github.com/gocityengineering/jira-metrics/jirametrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// run tests defined in datadir
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := flag.String("c", "/etc/jira-metrics/config.yaml", "configuration path")
	schemaPath := flag.String("s", "/etc/jira-metrics/schema.yaml", "schema path")

	flag.Parse()

	// fetch service account credentials
	username := os.Getenv("JIRA_METRICS_USERNAME")
	token := os.Getenv("JIRA_METRICS_API_TOKEN")

	os.Exit(realMain(*configPath, *schemaPath, username, token, false))
}

func realMain(configPath, schemaPath, username, token string, dryrun bool) int {

	// validate configuration file
	err := encodingutils.ValidateFile(configPath, schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config file %s is not valid against schema %s: %v\n", configPath, schemaPath, err)
		return 1
	}

	// parse configuration file
	var config = jirametrics.Config{
		Server:                                   "localhost",
		IntervalSeconds:                          300,
		Metrics:                                  []jirametrics.Metric{},
		Boards:                                   []jirametrics.Board{},
		Queues:                                   []jirametrics.Queue{},
		AverageSatisfactionScoresDayGauge:        prometheus.GaugeVec{},
		SatisfactionScoresDayGauge:               prometheus.GaugeVec{},
		OpenTicketsGauge:                         prometheus.GaugeVec{},
		NewTicketsDayGauge:                       prometheus.GaugeVec{},
		UnassignedTicketsGauge:                   prometheus.GaugeVec{},
		ResolvedTicketsDayGauge:                  prometheus.GaugeVec{},
		NewIncidents15mGauge:                     prometheus.GaugeVec{},
		NewTickets15mGauge:                       prometheus.GaugeVec{},
		OpenTicketsWithinSLAGauge:                prometheus.GaugeVec{},
		OpenTicketsOutsideSLAGauge:               prometheus.GaugeVec{},
		OpenTicketsOutsideSLAReporterCounter:     prometheus.CounterVec{},
		OpenTicketsOutsideSLAAssigneeCounter:     prometheus.CounterVec{},
		TicketsCreatedOutsideSLAReporterCounter:  prometheus.CounterVec{},
		TicketsCreatedOutsideSLAAssigneeCounter:  prometheus.CounterVec{},
		TicketsDoneGauge:                         prometheus.GaugeVec{},
		TicketsDoneBusinessAsUsualGauge:          prometheus.GaugeVec{},
		TicketsInPreparationGauge:                prometheus.GaugeVec{},
		TicketsInPreparationBusinessAsUsualGauge: prometheus.GaugeVec{},
		TicketsInProgressGauge:                   prometheus.GaugeVec{},
		TicketsInProgressBusinessAsUsualGauge:    prometheus.GaugeVec{},
		TicketsInPreparationReporterCounter:      prometheus.CounterVec{},
		TicketsInProgressReporterCounter:         prometheus.CounterVec{},
		TicketsDoneReporterCounter:               prometheus.CounterVec{},
	}
	err = jirametrics.ParseConfig(configPath, &config, dryrun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load config file %s: %v\n", configPath, err)
		return 2
	}

	// Jira API requests will fail unless we set InsecureSkipVerify to true
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// attempt minimal GET request
	err = jirametrics.MinimalGetRequest(config.Server, username, token, dryrun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't establish minimal connection to server: %v\n", err)
		return 3
	}

	go func() {
		for {
			err = jirametrics.UpdateMetrics(&config, username, token, dryrun)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't update metrics: %v\n", err)
			}
			if dryrun {
				return
			}
			time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
		}
	}()

	if dryrun {
		return 0
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
	return 0
}
