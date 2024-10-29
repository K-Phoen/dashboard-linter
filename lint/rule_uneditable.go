package lint

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewUneditableRule() *DashboardRuleFunc {
	return &DashboardRuleFunc{
		name:        "uneditable-dashboard",
		description: "Checks that the dashboard is not editable.",
		fn: func(d dashboard.Dashboard) DashboardRuleResults {
			r := DashboardRuleResults{}
			if d.Editable == nil || *d.Editable {
				r.AddFixableError(d, "is editable, it should be set to 'editable: false'", FixUneditableRule)
			}
			return r
		},
	}
}

func FixUneditableRule(d *dashboard.Dashboard) {
	d.Editable = cog.ToPtr(false)
}
