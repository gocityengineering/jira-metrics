package jirametrics

import (
	"fmt"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func processQueues(config *Config, username string, token string, dryrun bool) error {
	if dryrun {
		return nil
	}
	// reset metrics as required before iterating over queues
	config.NewIncidents15mGauge.Reset()
	config.NewTickets15mGauge.Reset()

	// process queues
	var jqlQuery string
	logLabel := "queue"
	for _, queue := range config.Queues {
		project := queue.Name
		safeName := queue.SafeName

		if queue.TimeZone == "" {
			queue.TimeZone = "UTC"
		}

		// satisfaction scores
		sum := 0
		divisor := 0
		for score := 1; score < 6; score++ {
			jqlQuery = fmt.Sprintf("project = %s and resolutiondate >= -1d and satisfaction is not null and satisfaction = %d", project, score)
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
				config.SatisfactionScoresDayGauge.With(prometheus.Labels{"queue": safeName, "score": strconv.Itoa(score)}).Set(float64(result.Total))
			}
			log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)
			// prepare avg
			sum += result.Total * score
			divisor += result.Total
		}
		if divisor > 0 && dryrun == false {
			config.AverageSatisfactionScoresDayGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(float64(sum) / float64(divisor)))
		}

		// new tickets
		jqlQuery = fmt.Sprintf("project = %s and created >= -1d", project)
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
			config.NewTicketsDayGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// resolved tickets
		jqlQuery = fmt.Sprintf("project = %s and status changed to resolved after -1d", project)
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
			config.ResolvedTicketsDayGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// unassigned tickets
		jqlQuery = fmt.Sprintf("project = %s and resolution is empty and assignee is empty", project)
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
			config.UnassignedTicketsGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// open tickets
		jqlQuery = fmt.Sprintf("project = %s and resolution is empty", project)
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
			config.OpenTicketsGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// new incidents (last 15m)
		jqlQuery = fmt.Sprintf("project = %s and created >= -15m and issuetype in (incident)", project)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)

		if err != nil {
			return err
		}
		if !dryrun {
			// flag all as out of hours if one is
			selfListString := ""
			summaryListString := ""

			for index, obj := range result.Issues {
				if index > 0 {
					selfListString = selfListString + ","
					summaryListString = summaryListString + ","
				}
				selfListString = selfListString + obj.Self
				summaryListString = summaryListString + obj.Fields.Summary
			}
			config.NewIncidents15mGauge.
				With(prometheus.Labels{"queue": safeName, "self": selfListString, "summary": summaryListString}).
				Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// new tickets (with out-of-hours flag for on-call use; last 15m)
		jqlQuery = fmt.Sprintf("project = %s and created >= -15m", project)
		result, err = SimpleQuery(
			config,
			username,
			token,
			dryrun,
			jqlQuery)

		if err != nil {
			return err
		}
		if !dryrun {
			// flag all as out of hours if one is
			outOfHoursString := "false"
			selfListString := ""
			summaryListString := ""

			for index, obj := range result.Issues {
				if index > 0 {
					selfListString = selfListString + ","
					summaryListString = summaryListString + ","
				}
				selfListString = selfListString + obj.Self
				summaryListString = summaryListString + obj.Fields.Summary
				created := obj.Fields.Created
				createdUnix, err := parseTimestamp(created)
				if err != nil {
					return err
				}
				outOfHours, err := dateStampOutOfHours(createdUnix, queue.TimeZone)
				if err != nil {
					return err
				}
				if outOfHours {
					outOfHoursString = "true"
				}
			}
			config.NewTickets15mGauge.
				With(prometheus.Labels{"queue": safeName, "self": selfListString, "summary": summaryListString, "outofhours": outOfHoursString}).
				Set(float64(result.Total))
		}
		log.Printf("[%s] %s / out of hours %d\n", logLabel, jqlQuery, result.Total)

		// tickets whose "Time to resolution" SLA is in breach
		jqlQuery = fmt.Sprintf("project='%s' AND resolution is empty AND 'Time to resolution'=breached()", project)
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
			config.OpenTicketsOutsideSLAGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// tickets whose "Time to resolution" SLA isn't in breach
		jqlQuery = fmt.Sprintf("project='%s' AND resolution is empty AND 'Time to resolution' != breached()", project)
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
			config.OpenTicketsWithinSLAGauge.With(prometheus.Labels{"queue": safeName}).Set(float64(result.Total))
		}
		log.Printf("[%s] %s %d\n", logLabel, jqlQuery, result.Total)

		// now use in-breach set to create counters showing reporters/assignees
		jqlQuery = fmt.Sprintf("project='%s' AND resolution is empty AND 'Time to resolution'=breached()", project)
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
				config.OpenTicketsOutsideSLAReporterCounter.With(prometheus.Labels{"queue": safeName, "reporter": reporter}).Inc()

				assignee := obj.Fields.Assignee.DisplayName
				config.OpenTicketsOutsideSLAAssigneeCounter.With(prometheus.Labels{"queue": safeName, "assignee": assignee}).Inc()
			}
		}
		log.Printf("[%s] %s / reporter / assignee %d\n", logLabel, jqlQuery, result.Total)
	}

	return nil
}
