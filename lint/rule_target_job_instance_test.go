package lint

import (
	"fmt"
	"testing"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func testTargetRequiredMatcherRule(t *testing.T, matcher string) {
	var linter *TargetRuleFunc

	switch matcher {
	case "job":
		linter = NewTargetJobRule()
	case "instance":
		linter = NewTargetInstanceRule()
	default:
		t.Errorf("No concrete target required matcher rule for '%s", matcher)
		return
	}

	for _, tc := range []struct {
		result Result
		target cogvariants.Dataquery
	}{
		// Happy path
		{
			result: ResultSuccess,
			target: prometheus.Dataquery{
				Expr: fmt.Sprintf(`sum(rate(foo{%s=~"$%s"}[5m]))`, matcher, matcher),
			},
		},
		// Also happy when the promql is invalid
		{
			result: ResultSuccess,
			target: prometheus.Dataquery{
				Expr: `foo(bar.baz))`,
			},
		},
		// Missing matcher
		{
			result: Result{
				Severity: Error,
				Message:  fmt.Sprintf("Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(rate(foo[5m]))': %s selector not found", matcher),
			},
			target: prometheus.Dataquery{
				Expr: `sum(rate(foo[5m]))`,
			},
		},
		// Not a regex matcher
		{
			result: Result{
				Severity: Error,
				Message:  fmt.Sprintf("Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(rate(foo{%s=\"$%s\"}[5m]))': %s selector is =, not =~", matcher, matcher, matcher),
			},
			target: prometheus.Dataquery{
				Expr: fmt.Sprintf(`sum(rate(foo{%s="$%s"}[5m]))`, matcher, matcher),
			},
		},
		// Wrong template variable
		{
			result: Result{
				Severity: Error,
				Message:  fmt.Sprintf("Dashboard 'dashboard', panel 'panel', target idx '0' invalid PromQL query 'sum(rate(foo{%s=~\"$foo\"}[5m]))': %s selector is $foo, not $%s", matcher, matcher, matcher),
			},
			target: prometheus.Dataquery{
				Expr: fmt.Sprintf(`sum(rate(foo{%s=~"$foo"}[5m]))`, matcher),
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
				{
					Panel: &dashboard.Panel{
						Title:   toPtr("panel"),
						Type:    "singlestat",
						Targets: []cogvariants.Dataquery{tc.target},
					},
				},
			},
		}

		testRule(t, linter, dash, tc.result)
	}
}

func TestTargetJobInstanceRule(t *testing.T) {
	testTargetRequiredMatcherRule(t, "job")
	testTargetRequiredMatcherRule(t, "instance")
}
