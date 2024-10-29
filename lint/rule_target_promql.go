package lint

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/prometheus/prometheus/promql/parser"
)

// panelHasQueries returns true is the panel has queries we should try and
// validate.  We allow-list panels here to prevent false positives with
// new panel types we don't understand.
func panelHasQueries(p *dashboard.Panel) bool {
	types := []string{panelTypeSingleStat, panelTypeGauge, panelTypeTimeTable, "stat", "state-timeline", panelTypeTimeSeries}
	for _, t := range types {
		if p.Type == t {
			return true
		}
	}
	return false
}

// parsePromQL returns the parsed PromQL statement from a panel,
// replacing eg [$__rate_interval] with [5m] so queries parse correctly.
// We also replace various other Grafana global variables.
func parsePromQL(expr string, variables []dashboard.VariableModel) (parser.Expr, error) {
	expr, err := expandVariables(expr, variables)
	if err != nil {
		return nil, fmt.Errorf("could not expand variables: %w", err)
	}
	return parser.ParseExpr(expr)
}

// NewTargetPromQLRule builds a lint rule for panels with Prometheus queries which checks:
// - the query is valid PromQL
// - the query contains two matchers within every selector - `{job=~"$job", instance=~"$instance"}`
// - the query is not empty
// - if the query references another panel then make sure that panel exists
func NewTargetPromQLRule() *TargetRuleFunc {
	return &TargetRuleFunc{
		name:        "target-promql-rule",
		description: "Checks that each target uses a valid PromQL query.",
		fn: func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel, t Target) TargetRuleResults {
			r := TargetRuleResults{}

			// The panel is a row
			if p.RowPanel != nil {
				return r
			}

			if t := getTemplateDatasource(d); t == nil || *t.Query.String != Prometheus {
				// Missing template datasources is a separate rule.
				return r
			}

			if !panelHasQueries(p.Panel) {
				return r
			}

			promQuery, ok := t.Original.(prometheus.Dataquery)
			if !ok {
				return r
			}

			// If panel does not contain an expression then check if it references another panel and it exists
			/*
					TODO(kgz): there's no PanelId in the foundation SDK :(
				if promQuery.Expr == nil || len(*promQuery.Expr) == 0 {
					if t.PanelId > 0 {
						for _, p1 := range d.Panels {
							if p1.Id == t.PanelId {
								return r
							}
						}
						r.AddError(d, p, t, "Invalid panel reference in target")
					}
				}
			*/

			if _, err := parsePromQL(promQuery.Expr, d.Templating.List); err != nil {
				r.AddError(d, p, t, fmt.Sprintf("invalid PromQL query '%s': %v", promQuery.Expr, err))
			}

			return r
		},
	}
}
