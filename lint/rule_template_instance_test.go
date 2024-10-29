package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestInstanceTemplate(t *testing.T) {
	linter := NewTemplateInstanceRule()

	for _, tc := range []struct {
		result    Result
		dashboard dashboard.Dashboard
	}{
		// Non-promtheus dashboards shouldn't fail.
		{
			result: ResultSuccess,
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
			},
		},
		// Missing instance templates.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' is missing the instance template",
			},
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
					},
				},
			},
		},
		// What success looks like.
		{
			result: ResultSuccess,
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
		testRule(t, linter, tc.dashboard, tc.result)
	}
}
