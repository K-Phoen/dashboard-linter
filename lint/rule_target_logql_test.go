package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func TestTargetLogQLRule(t *testing.T) {
	linter := NewTargetLogQLRule()

	for _, tc := range []struct {
		result Result
		panel  dashboard.Panel
	}{
		// Don't fail non-Loki panels.
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title:      toPtr("panel"),
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo"), Type: toPtr("prometheus")},
				Targets: []variants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[5m]))`,
					},
				},
			},
		},
		// Valid LogQL query
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job="mysql"}[5m]))`,
					},
				},
			},
		},
		// Invalid LogQL query
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' invalid LogQL query 'sum(rate({job="mysql"[5m]))': parse error at line 0, col 22: syntax error: unexpected RANGE, expecting } or ,`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job="mysql"[5m]))`,
					},
				},
			},
		},
		// Valid LogQL query with $__auto
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job="mysql"}[$__auto]))`,
					},
				},
			},
		},
		// Valid complex LogQL query
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum by (host) (rate({job="mysql"} |= "error" != "timeout" | json | duration > 10s [5m]))`,
					},
				},
			},
		},
		// Invalid complex LogQL query
		{
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'dashboard', panel 'panel', target idx '0' invalid LogQL query 'sum by (host) (rate({job="mysql"} |= "error" != "timeout" | json | duration > 10s [5m])))': parse error at line 1, col 89: syntax error: unexpected )`,
			},
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum by (host) (rate({job="mysql"} |= "error" != "timeout" | json | duration > 10s [5m])))`,
					},
				},
			},
		},
		// LogQL query with line_format
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `{job="mysql"} | json | line_format "{{.timestamp}} {{.message}}"`,
					},
				},
			},
		},
		// LogQL query with unwrap
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Title: toPtr("panel"),
				Type:  "singlestat",
				Targets: []variants.Dataquery{
					loki.Dataquery{
						Expr: `sum(rate({job="mysql"} | unwrap duration [5m]))`,
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
						Query: &dashboard.StringOrMap{String: toPtr("loki")},
					},
				},
			},
			Panels: []dashboard.PanelOrRowPanel{{Panel: &tc.panel}},
		}
		testRule(t, linter, dash, tc.result)
	}
}
