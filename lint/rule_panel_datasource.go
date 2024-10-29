package lint

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewPanelDatasourceRule() *PanelRuleFunc {
	return &PanelRuleFunc{
		name:        "panel-datasource-rule",
		description: "Checks that each panel uses the templated datasource.",
		fn: func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel) PanelRuleResults {
			r := PanelRuleResults{}

			// The panel is a row
			if p.RowPanel != nil {
				return r
			}

			switch p.Panel.Type {
			case panelTypeSingleStat, panelTypeGraph, panelTypeTimeTable, panelTypeTimeSeries:
				// That a templated datasource exists, is the responsibility of another rule.
				templatedDs := getTemplateByType(d, "datasource")
				availableDsUids := make(map[string]struct{}, len(templatedDs)*2)
				for _, tds := range templatedDs {
					availableDsUids[fmt.Sprintf("$%s", tds.Name)] = struct{}{}
					availableDsUids[fmt.Sprintf("${%s}", tds.Name)] = struct{}{}
				}

				datasourceUid := ""
				if p.Panel.Datasource != nil && p.Panel.Datasource.Uid != nil {
					datasourceUid = *p.Panel.Datasource.Uid
				}

				_, ok := availableDsUids[datasourceUid]
				if !ok {
					r.AddError(d, p, fmt.Sprintf("does not use a templated datasource, uses '%s'", datasourceUid))
				}
			}

			return r
		},
	}
}
