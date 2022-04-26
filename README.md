# jira-metrics

This service communicates directly with Jira.

All teams are welcome to add metrics to `chart/templates/configmap.yaml`. Here is an excerpt:

```
data:
  config.yaml: |
    server: acme.atlassian.net
    metrics:
      -
        query: 'project = ACME AND status changed to "waiting for support" after -1d'
        metric_name: helpdesk_tickets_new_1d
        metric_type: gauge
        metric_help: 'new helpdesk tickets'
```

The `server` will be the same for all queries. If your metrics are not coming through, the user specified in the secret may require additional access to your tickets.

You can enter any JQL query in the field `query`, but the result's `Total` property (i.e. the number of results) will be submitted to Prometheus. It should make sense in relation to a given workday; often that will mean adding `after -1d` to your JQL expression, but at other times what you are tracking may be something like the total number of open tickets, in which case no time bracket is required.

The `metric_name` is the string you will query in Prometheus/Grafana. The usual recommendations apply. Be mindful that all Prometheus users will see all metrics, so it helps to use names that clearly identify a given metric as yours.

The `metric_type` could be `counter`, `gauge` or `histogram`. For business data, you may find `gauge` and `histogram` more useful than the `counter` collector, which works best with high-frequency inputs such as server response codes.

The `metric_help` field describes your metric in human-friendly terms.

