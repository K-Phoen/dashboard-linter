package lint

import (
	"fmt"
	"testing"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func TestTargetPromQLRule(t *testing.T) {
	linter := NewTargetPromQLRule()

	for i, tc := range []struct {
		result []Result
		panel  dashboard.Panel
	}{
		// Don't fail non-prometheus panels.
		{
			result: []Result{ResultSuccess},
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
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[5m]))`,
					},
				},
			},
		},
		// Invalid query
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[5m]))`,
					},
				},
			},
		},
		// Timeseries support
		{
			result: []Result{{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'foo(bar.baz)': 1:8: parse error: unexpected character: '.'",
			}},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "timeseries",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `foo(bar.baz)`,
					},
				},
			},
		},
		// Variable substitutions
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[$__rate_interval])) * $__range_s`,
					},
				},
			},
		},
		// Variable substitutions with ${...}
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[$__rate_interval])) * ${__range_s}`,
					},
				},
			},
		},
		// Variable substitutions inside by clause
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum by(${variable:csv}) (rate(foo[$__rate_interval])) * $__range_s`,
					},
				},
			},
		},
		// Template variables substitutions
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum (rate(foo[$interval:$resolution]))`,
					},
				},
			},
		},
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `increase(foo{}[$sampling])`,
					},
				},
			},
		},
		// Empty PromQL expression
		{
			result: []Result{{
				Severity: Error,
				Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query '': unknown position: parse error: no expression found in input",
			}},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: ``,
					},
				},
			},
		},
		// Reference another panel that does not exist
		/*
			{
				result: []Result{
					{
						Severity: Error,
						Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' Invalid panel reference in target",
					},
					{
						Severity: Error,
						Message:  "Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query '': unknown position: parse error: no expression found in input",
					},
				},
				panel: dashboard.Panel{
					Id:    toPtr[uint32](1),
					Title: toPtr("panel"),
					Type:  "singlestat",
					Targets: []cogvariants.Dataquery{
						prometheus.Dataquery{
							PanelId: 2, // PanelId does not exist in Foundation SDK (not in the schema)
						},
					},
				},
			},
		*/
	} {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			d := dashboard.Dashboard{
				Title: toPtr("dashboard"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Type: "interval",
							Name: "interval",
							Options: []dashboard.VariableOption{
								{Value: dashboard.StringOrArrayOfString{String: toPtr("1h")}},
							},
						},
						{
							Type: "interval",
							Name: "sampling",
							Current: &dashboard.VariableOption{
								Value: dashboard.StringOrArrayOfString{String: toPtr("$__auto_interval_sampling")},
							},
						},
						{
							Type: "resolution",
							Name: "resolution",
							Options: []dashboard.VariableOption{
								{Value: dashboard.StringOrArrayOfString{String: toPtr("1h")}},
								{Value: dashboard.StringOrArrayOfString{String: toPtr("1h")}},
							},
						},
					},
				},
				Panels: []dashboard.PanelOrRowPanel{
					{Panel: &tc.panel},
				},
			}

			testMultiResultRule(t, linter, d, tc.result)
		})
	}
}
