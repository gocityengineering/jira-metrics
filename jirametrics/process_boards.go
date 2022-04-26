package jirametrics

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func processBoards(config *Config, username string, token string, dryrun bool) error {
	if dryrun {
		return nil
	}
	var jqlQuery string
	logLabel := "board"
	for _, board := range config.Boards {
		project := board.Name
		safeName := board.SafeName
		iterationLengthWeeks := board.IterationLengthWeeks
		businessAsUsualEpics := board.BusinessAsUsualEpics
		inProgressLabels := board.InProgressLabels
		inPreparationLabels := board.InPreparationLabels

		// GAUGES

		// TicketsDoneGauge

		jqlQuery = fmt.Sprintf("project = %s and status changed to 'done' after -%dw", project, iterationLengthWeeks)
		result, err := SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsDoneGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// TicketsDoneBusinessAsUsualGauge
		jqlQuery = fmt.Sprintf("project = %s and status changed to 'done' after -%dw and 'Epic Link' in (%s)", project, iterationLengthWeeks, businessAsUsualEpics)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsDoneBusinessAsUsualGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// TicketsInPreparationGauge
		jqlQuery = fmt.Sprintf("project = %s and status in (%s)", project, inPreparationLabels)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsInPreparationGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// TicketsInPreparationBusinessAsUsualGauge
		jqlQuery = fmt.Sprintf("project = %s and status in (%s) and 'Epic Link' in (%s)", project, inPreparationLabels, businessAsUsualEpics)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsInPreparationBusinessAsUsualGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// TicketsInProgressGauge
		jqlQuery = fmt.Sprintf("project = %s and status in (%s)", project, inProgressLabels)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsInProgressGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// TicketsInProgressBusinessAsUsualGauge
		jqlQuery = fmt.Sprintf("project = %s and status in (%s) and 'Epic Link' in (%s)", project, inProgressLabels, businessAsUsualEpics)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			config.TicketsInProgressBusinessAsUsualGauge.With(prometheus.Labels{"board": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// COUNTERS

		// TicketDoneReporterCounter
		jqlQuery = fmt.Sprintf("project = %s and resolved > -%dw", project, iterationLengthWeeks)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			for _, obj := range result.Issues {
				reporter := obj.Fields.Reporter.DisplayName
				config.TicketsDoneReporterCounter.With(prometheus.Labels{"board": safeName, "reporter": reporter}).Inc()
			}
		}

		// TicketsInPreparationReporterCounter
		jqlQuery = fmt.Sprintf("project = %s and status in (%s)", project, inPreparationLabels)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			for _, obj := range result.Issues {
				reporter := obj.Fields.Reporter.DisplayName
				config.TicketsInPreparationReporterCounter.With(prometheus.Labels{"board": safeName, "reporter": reporter}).Inc()
			}
		}

		// TicketsInProgressReporterCounter
		jqlQuery = fmt.Sprintf("project = %s and status in (%s)", project, inProgressLabels)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)
		if err != nil {
			return err
		}
		if dryrun == false {
			for _, obj := range result.Issues {
				reporter := obj.Fields.Reporter.DisplayName
				config.TicketsInProgressReporterCounter.With(prometheus.Labels{"board": safeName, "reporter": reporter}).Inc()
			}
		}

	}

	return nil
}
