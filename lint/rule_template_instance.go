package lint

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewTemplateInstanceRule() *DashboardRuleFunc {
	return &DashboardRuleFunc{
		name:        "template-instance-rule",
		description: "Checks that the dashboard has a templated instance.",
		fn: func(d dashboard.Dashboard) DashboardRuleResults {
			r := DashboardRuleResults{}

			template := getTemplateDatasource(d)
			if template == nil || *template.Query.String != Prometheus {
				return r
			}

			checkTemplate(d, "instance", &r)
			return r
		},
	}
}
