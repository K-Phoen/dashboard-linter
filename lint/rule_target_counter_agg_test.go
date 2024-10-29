package lint

import (
	"testing"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func TestTargetCounterAggRule(t *testing.T) {
	linter := NewTargetCounterAggRule()

	for _, tc := range []struct {
		result Result
		panel  dashboard.Panel
	}{
		// Non aggregated counter fails
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' counter metric 'something_total' is not aggregated with rate, irate, or increase",
			},
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `something_total`,
					},
				},
			},
		},
		// Weird matrix selector without an aggregator
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' counter metric 'something_total' is not aggregated with rate, irate, or increase",
			},
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `something_total[$__rate_interval]`,
					},
				},
			},
		},
		// Single aggregated counter is good
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `increase(something_total[$__rate_interval])`,
					},
				},
			},
		},
		// Sanity check for multiple counters in one query, with the first one failing
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' counter metric 'something_total' is not aggregated with rate, irate, or increase",
			},
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `something_total / rate(somethingelse_total[$__rate_interval])`,
					},
				},
			},
		},
		// Sanity check for multiple counters in one query, with the second one failing
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' counter metric 'somethingelse_total' is not aggregated with rate, irate, or increase",
			},
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `rate(something_total[$__rate_interval]) / somethingelse_total`,
					},
				},
			},
		},
	} {
		panels := []dashboard.PanelOrRowPanel{
			{Panel: &tc.panel},
		}
		testRule(t, linter, dashboard.Dashboard{Title: toPtr("dashboard"), Panels: panels}, tc.result)
	}
}
