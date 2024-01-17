package lint

import (
	"testing"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func TestTargetRateIntervalRule(t *testing.T) {
	linter := NewTargetRateIntervalRule()

	for _, tc := range []struct {
		result Result
		panel  dashboard.Panel
	}{
		// Don't fail non-prometheus panels.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[5m]))`,
					},
				},
			},
		},
		// This is what a valid panel looks like.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo{job=~"$job",instance=~"$instance"}[$__rate_interval]))`,
					},
				},
			},
		},
		// This is what a valid panel looks like.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo{job=~"$job",instance=~"$instance"}[$__rate_interval]))/sum(rate(bar{job=~"$job",instance=~"$instance"}[$__rate_interval]))`,
					},
				},
			},
		},
		// Invalid query
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(rate(foo{job=~"$job",instance=~"$instance"}[5m]))': should use $__rate_interval`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo{job=~"$job",instance=~"$instance"}[5m]))`,
					},
				},
			},
		},
		// Timeseries support
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(rate(foo{job=~"$job",instance=~"$instance"}[5m]))': should use $__rate_interval`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "timeseries",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo{job=~"$job",instance=~"$instance"}[5m]))`,
					},
				},
			},
		},
		// Non-rate functions should not make the linter fail
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(increase(foo{job=~"$job",instance=~"$instance"}[$__range]))`,
					},
				},
			},
		},
		// irate should be checked too
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(irate(foo{job=~"$job",instance=~"$instance"}[$__interval]))': should use $__rate_interval`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(irate(foo{job=~"$job",instance=~"$instance"}[$__interval]))`,
					},
				},
			},
		},
	} {
		dash := dashboard.Dashboard{
			Title: toPtr("dashboard"),
			Templating: dashboard.DashboardDashboardTemplating{
				List: []dashboard.VariableModel{
					{
						Type:  "datasource",
						Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
					},
				},
			},
			Panels: []dashboard.PanelOrRowPanel{
				{Panel: &tc.panel},
			},
		}

		testRule(t, linter, dash, tc.result)
	}
}
