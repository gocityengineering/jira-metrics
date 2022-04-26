package jirametrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Server          string   `json:"server"`
	Metrics         []Metric `json:"metrics"`
	Boards          []Board  `json:"boards"`
	Queues          []Queue  `json:"queues"`
	IntervalSeconds int      `json:"interval_seconds"`
	// queues
	AverageSatisfactionScoresDayGauge       prometheus.GaugeVec
	SatisfactionScoresDayGauge              prometheus.GaugeVec
	OpenTicketsGauge                        prometheus.GaugeVec
	NewTicketsDayGauge                      prometheus.GaugeVec
	UnassignedTicketsGauge                  prometheus.GaugeVec
	ResolvedTicketsDayGauge                 prometheus.GaugeVec
	NewIncidents15mGauge                    prometheus.GaugeVec
	NewTickets15mGauge                      prometheus.GaugeVec
	OpenTicketsWithinSLAGauge               prometheus.GaugeVec
	OpenTicketsOutsideSLAGauge              prometheus.GaugeVec
	OpenTicketsOutsideSLAReporterCounter    prometheus.CounterVec
	OpenTicketsOutsideSLAAssigneeCounter    prometheus.CounterVec
	TicketsCreatedOutsideSLAReporterCounter prometheus.CounterVec
	TicketsCreatedOutsideSLAAssigneeCounter prometheus.CounterVec
	// boards
	TicketsDoneGauge                         prometheus.GaugeVec
	TicketsDoneBusinessAsUsualGauge          prometheus.GaugeVec
	TicketsInPreparationGauge                prometheus.GaugeVec
	TicketsInPreparationBusinessAsUsualGauge prometheus.GaugeVec
	TicketsInProgressGauge                   prometheus.GaugeVec
	TicketsInProgressBusinessAsUsualGauge    prometheus.GaugeVec
	TicketsInPreparationReporterCounter      prometheus.CounterVec
	TicketsInProgressReporterCounter         prometheus.CounterVec
	TicketsDoneReporterCounter               prometheus.CounterVec
}

type Metric struct {
	Query      string `json:"query"`
	MetricName string `json:"metric_name"`
	MetricType string `json:"metric_type"`
	MetricHelp string `json:"metric_help"`
	Counter    prometheus.Counter
	Gauge      prometheus.Gauge
	Histogram  prometheus.Histogram
}

type Board struct {
	Name                 string `json:"name"`
	InProgressLabels     string `json:"in_progress_labels"`
	InPreparationLabels  string `json:"in_preparation_labels"`
	BusinessAsUsualEpics string `json:"business_as_usual_epics"`
	SafeName             string `json:"safe_name"`
	IterationLengthWeeks int    `json:"iteration_length_weeks"`
}

type Queue struct {
	Name     string `json:"name"`
	TimeZone string `json:"time_zone"`
	SafeName string `json:"safe_name"`
}

type Data struct {
	Jql string `json:"jql"`
}

type Result struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

type Issue struct {
	Fields Field  `json:"fields"`
	Self   string `json:"key"` // this is the ticket number
}

type Field struct {
	Assignee       Person `json:"assignee"`
	Created        string `json:"created"`
	Reporter       Person `json:"reporter"`
	ResolutionDate string `json:"resolutiondate"`
	Summary        string `json:"summary"`
}

type Person struct {
	DisplayName string `json:"displayName"`
}
