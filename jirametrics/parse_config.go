package jirametrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func ParseConfig(configPath string, config *Config, dryrun bool) error {
	byteArray, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("can't read configuration file %s: %v", configPath, err)
	}

	jsonArray, err := yaml.YAMLToJSON(byteArray)
	if err != nil {
		return fmt.Errorf("can't convert configuration file %s to JSON: %v", configPath, err)
	}

	err = json.Unmarshal(jsonArray, config)

	// set up appropriate Collector
	for index, metric := range config.Metrics {
		log.Printf("Registering %s\n", metric.MetricName)
		switch metric.MetricType {
		case "counter":
			config.Metrics[index].Counter = promauto.NewCounter(prometheus.CounterOpts{
				Subsystem: "jira_metrics",
				Name:      metric.MetricName,
				Help:      metric.MetricHelp,
			})
		case "gauge":
			config.Metrics[index].Gauge = promauto.NewGauge(prometheus.GaugeOpts{
				Subsystem: "jira_metrics",
				Name:      metric.MetricName,
				Help:      metric.MetricHelp,
			})
		case "histogram":
			config.Metrics[index].Histogram = promauto.NewHistogram(prometheus.HistogramOpts{
				Subsystem: "jira_metrics",
				Name:      metric.MetricName,
				Help:      metric.MetricHelp,
			})
		default:
			return fmt.Errorf("metric %s missing recognised type", metric.MetricName)
		}
	}

	for index, queue := range config.Queues {
		log.Printf("Registering queue %s\n", queue.Name)

		// prepare safe name for label use
		safe := strings.ToLower(queue.Name)
		regex, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return err
		}
		safe = regex.ReplaceAllString(safe, "")

		config.Queues[index].SafeName = safe
	}

	for index, board := range config.Boards {
		log.Printf("Registering board %s\n", board.Name)

		// prepare safe name for label use
		safe := strings.ToLower(board.Name)
		regex, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return err
		}
		safe = regex.ReplaceAllString(safe, "")

		config.Boards[index].SafeName = safe
	}

	if dryrun {
		return nil
	}

	// queues
	config.SatisfactionScoresDayGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_satisfaction_scores_1d",
		Help: "gauge for CSAT score",
	},
		[]string{
			"queue",
			"score",
		})
	prometheus.MustRegister(config.SatisfactionScoresDayGauge)

	config.AverageSatisfactionScoresDayGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_average_satisfaction_scores_1d",
		Help: "gauge for average CSAT score",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.AverageSatisfactionScoresDayGauge)

	config.NewTicketsDayGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_new_tickets_1d",
		Help: "gauge for new queue tickets",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.NewTicketsDayGauge)

	config.ResolvedTicketsDayGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_resolved_tickets_1d",
		Help: "gauge for resolved queue tickets",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.ResolvedTicketsDayGauge)

	config.OpenTicketsGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_open_tickets",
		Help: "gauge for open queue tickets",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.OpenTicketsGauge)

	config.UnassignedTicketsGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_unassigned_tickets",
		Help: "gauge for open queue tickets",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.UnassignedTicketsGauge)

	config.NewIncidents15mGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_new_incidents_15m",
		Help: "gauge for new incidents",
	},
		[]string{
			"queue",
			"self",
			"summary",
		})
	prometheus.MustRegister(config.NewIncidents15mGauge)

	config.NewTickets15mGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_new_tickets_15m",
		Help: "gauge for new tickets with out-of-hours flag",
	},
		[]string{
			"queue",
			"self",
			"summary",
			"outofhours",
		})
	prometheus.MustRegister(config.NewTickets15mGauge)

	// SLA-focused metrics
	config.OpenTicketsWithinSLAGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_open_tickets_within_sla",
		Help: "gauge for open tickets within SLA",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.OpenTicketsWithinSLAGauge)

	config.OpenTicketsOutsideSLAGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_queue_open_tickets_outside_sla",
		Help: "gauge for open tickets outside SLA",
	},
		[]string{
			"queue",
		})
	prometheus.MustRegister(config.OpenTicketsOutsideSLAGauge)

	config.OpenTicketsOutsideSLAReporterCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_queue_open_tickets_outside_sla_reporter_total",
		Help: "SLA breach counter by reporter",
	},
		[]string{
			"queue",
			"reporter",
		})
	prometheus.MustRegister(config.OpenTicketsOutsideSLAReporterCounter)

	config.OpenTicketsOutsideSLAAssigneeCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_queue_open_tickets_outside_sla_assignee_total",
		Help: "SLA breach counter by assignee",
	},
		[]string{
			"queue",
			"assignee",
		})
	prometheus.MustRegister(config.OpenTicketsOutsideSLAAssigneeCounter)

	config.TicketsCreatedOutsideSLAReporterCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_queue_tickets_created_outside_sla_reporter_total",
		Help: "SLA breach counter (created) by reporter",
	},
		[]string{
			"queue",
			"reporter",
		})
	prometheus.MustRegister(config.TicketsCreatedOutsideSLAReporterCounter)

	config.TicketsCreatedOutsideSLAAssigneeCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_queue_tickets_created_outside_sla_assignee_total",
		Help: "SLA breach counter (created) by assignee",
	},
		[]string{
			"queue",
			"assignee",
		})
	prometheus.MustRegister(config.TicketsCreatedOutsideSLAAssigneeCounter)

	// boards
	config.TicketsDoneGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_done",
		Help: "Board tickets done",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsDoneGauge)

	config.TicketsInPreparationGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_in_preparation",
		Help: "Board tickets in preparation",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsInPreparationGauge)

	config.TicketsInProgressGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_in_progress",
		Help: "Board tickets in progress",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsInProgressGauge)

	config.TicketsInPreparationBusinessAsUsualGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_in_preparation_business_as_usual",
		Help: "Board tickets in preparation business as usual",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsInPreparationBusinessAsUsualGauge)

	config.TicketsInProgressBusinessAsUsualGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_in_progress_business_as_usual",
		Help: "Board tickets in progress business as usual",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsInProgressBusinessAsUsualGauge)

	config.TicketsDoneBusinessAsUsualGauge = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jira_metrics_board_tickets_done_business_as_usual",
		Help: "Board tickets done business as usual",
	},
		[]string{
			"board",
		})
	prometheus.MustRegister(config.TicketsDoneBusinessAsUsualGauge)

	config.TicketsInPreparationReporterCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_board_tickets_in_preparation_reporter_total",
		Help: "Board tickets in preparation by reporter",
	},
		[]string{
			"board",
			"reporter",
		})
	prometheus.MustRegister(config.TicketsInPreparationReporterCounter)

	config.TicketsInProgressReporterCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_board_tickets_in_progress_reporter_total",
		Help: "Board tickets in progress by reporter",
	},
		[]string{
			"board",
			"reporter",
		})
	prometheus.MustRegister(config.TicketsInProgressReporterCounter)

	config.TicketsDoneReporterCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "jira_metrics_board_tickets_done_reporter_total",
		Help: "Board tickets done by reporter",
	},
		[]string{
			"board",
			"reporter",
		})
	prometheus.MustRegister(config.TicketsDoneReporterCounter)

	return nil
}
