package lint

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NewTemplateJobRule() *DashboardRuleFunc {
	return &DashboardRuleFunc{
		name:        "template-job-rule",
		description: "Checks that the dashboard has a templated job.",
		fn: func(d dashboard.Dashboard) DashboardRuleResults {
			r := DashboardRuleResults{}

			template := getTemplateDatasource(d)
			if template == nil || *template.Query.String != Prometheus {
				return r
			}

			checkTemplate(d, "job", &r)
			return r
		},
	}
}

func checkTemplate(d dashboard.Dashboard, name string, r *DashboardRuleResults) {
	t := getTemplate(d, name)
	if t == nil {
		r.AddError(d, fmt.Sprintf("is missing the %s template", name))
		return
	}

	// TODO: Adding the prometheus_datasource here is hacky. This check function also assumes that all template vars which it will
	// ever check are only prometheus queries, which may not always be the case.
	datasourceUid := ""
	if t.Datasource != nil && t.Datasource.Uid != nil {
		datasourceUid = *t.Datasource.Uid
	}

	if datasourceUid != "$datasource" && datasourceUid != "${datasource}" && datasourceUid != "$prometheus_datasource" && datasourceUid != "${prometheus_datasource}" {
		r.AddError(d, fmt.Sprintf("%s template should use datasource '$datasource', is currently '%s'", name, datasourceUid))
	}

	if t.Type != targetTypeQuery {
		r.AddError(d, fmt.Sprintf("%s template should be a Prometheus query, is currently '%s'", name, t.Type))
	}

	titleCaser := cases.Title(language.English)
	labelTitle := titleCaser.String(name)

	label := ""
	if t.Label != nil {
		label = *t.Label
	}

	if label != labelTitle {
		r.AddWarning(d, fmt.Sprintf("%s template should be a labeled '%s', is currently '%s'", name, labelTitle, label))
	}

	if t.Multi == nil || !*t.Multi {
		r.AddError(d, fmt.Sprintf("%s template should be a multi select", name))
	}

	if t.AllValue == nil || *t.AllValue != ".+" {
		allValue := ""
		if t.AllValue != nil {
			allValue = *t.AllValue
		}
		r.AddError(d, fmt.Sprintf("%s template allValue should be '.+', is currently '%s'", name, allValue))
	}
}
