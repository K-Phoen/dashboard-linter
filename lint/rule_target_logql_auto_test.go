package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

// TestTargetLogQLAutoRule tests the NewTargetLogQLAutoRule function to ensure
// that it correctly identifies LogQL queries that should use $__auto for range vectors.
func TestTargetLogQLAutoRule(t *testing.T) {
	linter := NewTargetLogQLAutoRule()

	for _, tc := range []struct {
		result Result
		panel  dashboard.Panel
	}{
		// Test case: Non-Loki panel should pass without errors.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo"), Type: toPtr("prometheus")},
				Targets: []variants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate({job=~"$job",instance=~"$instance"}[5m]))`,
					},
				},
			},
		},
		// Test case: Valid LogQL query using $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job=~"$job",instance=~"$instance"} [$__auto]))`,
					},
				},
			},
		},
		// Test case: Valid LogQL query using $__auto in a complex expression.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job=~"$job",instance=~"$instance"} [$__auto]))/sum(rate({job=~"$job",instance=~"$instance"} [$__auto]))`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job=~"$job",instance=~"$instance"}[5m]))`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query without $__auto in a timeseries panel.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "timeseries",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job=~"$job",instance=~"$instance"}[5m]))`,
					},
				},
			},
		},
		// Test case: Valid LogQL query with count_over_time and $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `count_over_time({job="mysql"} [$__auto])`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query with count_over_time without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `count_over_time({job="mysql"}[5m])`,
					},
				},
			},
		},
		// Test case: Valid LogQL query with bytes_rate and $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `bytes_rate({job="mysql"} [$__auto])`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query with bytes_rate without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `bytes_rate({job="mysql"}[5m])`,
					},
				},
			},
		},
		// Test case: Valid LogQL query with bytes_over_time and $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `bytes_over_time({job="mysql"} [$__auto])`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query with bytes_over_time without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `bytes_over_time({job="mysql"}[5m])`,
					},
				},
			},
		},
		// Test case: Valid LogQL query with sum_over_time and $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum_over_time({job="mysql"} |= "duration" | unwrap duration [$__auto])`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query with sum_over_time without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum_over_time({job="mysql"} |= "duration" | unwrap duration[5m])`,
					},
				},
			},
		},
		// Test case: Valid LogQL query with avg_over_time and $__auto.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `avg_over_time({job="mysql"} |= "duration" | unwrap duration [$__auto])`,
					},
				},
			},
		},
		// Test case: Invalid LogQL query with avg_over_time without $__auto.
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' LogQL query uses fixed duration: should use $__auto`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `avg_over_time({job="mysql"} |= "duration" | unwrap duration[5m])`,
					},
				},
			},
		},
		// Add similar tests for other unwrapped range aggregations...
	} {
		dash := dashboard.Dashboard{
			Title: toPtr("dashboard"),
			Templating: dashboard.DashboardDashboardTemplating{
				List: []dashboard.VariableModel{
					{
						Type:  "datasource",
						Query: &dashboard.StringOrMap{String: toPtr("loki")},
					},
				},
			},
			Panels: []dashboard.PanelOrRowPanel{{Panel: &tc.panel}},
		}
		testRule(t, linter, dash, tc.result)
	}
}
