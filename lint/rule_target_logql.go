package lint

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
)

func NewTargetLogQLRule() *TargetRuleFunc {
	return &TargetRuleFunc{
		name:        "target-logql-rule",
		description: "Checks that each target uses a valid LogQL query.",
		fn: func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel, t Target) TargetRuleResults {
			r := TargetRuleResults{}

			// The panel is a row
			if p.RowPanel != nil {
				return r
			}

			lokiQuery, ok := t.Original.(loki.Dataquery)
			if !ok {
				return r
			}

			// Skip hidden targets
			if lokiQuery.Hide != nil && *lokiQuery.Hide {
				return r
			}

			// Check if the datasource is Loki
			isLoki := false
			if templateDS := getTemplateDatasource(d); templateDS != nil && templateDS.Query.String != nil && *templateDS.Query.String == Loki {
				isLoki = true
			} else if lokiQuery.Datasource != nil && lokiQuery.Datasource.Type != nil && *lokiQuery.Datasource.Type == Loki {
				isLoki = true
			}

			// skip if the datasource is not Loki
			if !isLoki {
				return r
			}

			if !panelHasQueries(p.Panel) {
				return r
			}

			// If panel does not contain an expression then check if it references another panel and it exists
			if len(lokiQuery.Expr) == 0 {
				/*
					TODO(kgz): there's no PanelId in the foundation SDK :(
					if t.PanelId > 0 {
						for _, p1 := range d.Panels {
							if p1.Id == t.PanelId {
								return r
							}
						}
						r.AddError(d, p, t, "Invalid panel reference in target")
					}
				*/
				return r
			}

			// Parse the LogQL query
			_, err := parseLogQL(lokiQuery.Expr, d.Templating.List)
			if err != nil {
				r.AddError(d, p, t, fmt.Sprintf("invalid LogQL query '%s': %v", lokiQuery.Expr, err))
				return r
			}

			return r
		},
	}
}
