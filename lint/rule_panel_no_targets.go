package lint

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewPanelNoTargetsRule() *PanelRuleFunc {
	return &PanelRuleFunc{
		name:        "panel-no-targets-rule",
		description: "Checks that each panel has at least one target.",
		fn: func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel) PanelRuleResults {
			r := PanelRuleResults{}

			// The panel is a row
			if p.RowPanel != nil {
				return r
			}

			switch p.Panel.Type {
			case panelTypeStat, panelTypeSingleStat, panelTypeGraph, panelTypeTimeTable, panelTypeTimeSeries, panelTypeGauge:
				if p.Panel.Targets != nil {
					return r
				}

				r.AddError(d, p, "has no targets")
			}
			return r
		},
	}
}
