package lint

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewTemplateOnTimeRangeReloadRule() *DashboardRuleFunc {
	return &DashboardRuleFunc{
		name:        "template-on-time-change-reload-rule",
		description: "Checks that the dashboard template variables are configured to reload on time change.",
		fn: func(d dashboard.Dashboard) DashboardRuleResults {
			r := DashboardRuleResults{}

			for i, template := range d.Templating.List {
				if template.Type != targetTypeQuery {
					continue
				}

				if template.Refresh != nil && *template.Refresh != dashboard.VariableRefreshOnTimeRangeChanged {
					r.AddFixableError(d,
						fmt.Sprintf("templated datasource variable named '%s', should be set to be refreshed "+
							"'On Time Range Change (value 2)', is currently '%d'", template.Name, *template.Refresh),
						fixTemplateOnTimeRangeReloadRule(i))
				}
			}
			return r
		},
	}
}

func fixTemplateOnTimeRangeReloadRule(i int) func(*dashboard.Dashboard) {
	return func(d *dashboard.Dashboard) {
		refresh := dashboard.VariableRefreshOnTimeRangeChanged
		d.Templating.List[i].Refresh = &refresh
	}
}
