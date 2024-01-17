package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestJobTemplate(t *testing.T) {
	linter := NewTemplateJobRule()

	for _, tc := range []struct {
		name      string
		result    []Result
		dashboard dashboard.Dashboard
	}{
		{
			name:   "Non-promtheus dashboards shouldn't fail.",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
			},
		},
		{
			name: "Missing job template.",
			result: []Result{{
				Severity: Error,
				Message:  "Dashboard 'test' is missing the job template",
			}},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
					},
				},
			},
		},
		{
			name: "Wrong datasource.",
			result: []Result{
				{Severity: Error, Message: "Dashboard 'test' job template should use datasource '$datasource', is currently 'foo'"},
				{Severity: Error, Message: "Dashboard 'test' job template should be a Prometheus query, is currently ''"},
				{Severity: Warning, Message: "Dashboard 'test' job template should be a labeled 'Job', is currently ''"},
				{Severity: Error, Message: "Dashboard 'test' job template should be a multi select"},
				{Severity: Error, Message: "Dashboard 'test' job template allValue should be '.+', is currently ''"}},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Name:       "job",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
						},
					},
				},
			},
		},
		{
			name: "Wrong type.",
			result: []Result{
				{Severity: Error, Message: "Dashboard 'test' job template should be a Prometheus query, is currently 'bar'"},
				{Severity: Warning, Message: "Dashboard 'test' job template should be a labeled 'Job', is currently ''"},
				{Severity: Error, Message: "Dashboard 'test' job template should be a multi select"},
				{Severity: Error, Message: "Dashboard 'test' job template allValue should be '.+', is currently ''"}},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Name:       "job",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Type:       "bar",
						},
					},
				},
			},
		},
		{
			name: "Wrong job label.",
			result: []Result{
				{Severity: Warning, Message: "Dashboard 'test' job template should be a labeled 'Job', is currently 'bar'"},
				{Severity: Error, Message: "Dashboard 'test' job template should be a multi select"},
				{Severity: Error, Message: "Dashboard 'test' job template allValue should be '.+', is currently ''"}},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Name:       "job",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Type:       "query",
							Label:      toPtr("bar"),
						},
					},
				},
			},
		},
		{
			name:   "OK",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Name:       "job",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Type:       "query",
							Label:      toPtr("Job"),
							Multi:      toPtr(true),
							AllValue:   toPtr(".+"),
						},
						{
							Name:       "instance",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("${datasource}")},
							Type:       "query",
							Label:      toPtr("Instance"),
							Multi:      toPtr(true),
							AllValue:   toPtr(".+"),
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testMultiResultRule(t, linter, tc.dashboard, tc.result)
		})
	}
}
