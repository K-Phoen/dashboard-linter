package lint

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewPanelTitleDescriptionRule() *PanelRuleFunc {
	return &PanelRuleFunc{
		name:        "panel-title-description-rule",
		description: "Checks that each panel has a title and description.",
		fn: func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel) PanelRuleResults {
			r := PanelRuleResults{}

			// The panel is a row
			if p.RowPanel != nil {
				return r
			}

			switch p.Panel.Type {
			case panelTypeStat, panelTypeSingleStat, panelTypeGraph, panelTypeTimeTable, panelTypeTimeSeries, panelTypeGauge:
				if p.Panel.Title == nil || len(*p.Panel.Title) == 0 {
					r.AddError(d, p, "has missing title")
				}

				if p.Panel.Description == nil || len(*p.Panel.Description) == 0 {
					r.AddError(d, p, "has missing description")
				}
			}
			return r
		},
	}
}
