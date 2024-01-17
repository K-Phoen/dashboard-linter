package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestTemplateLabelPromQLRule(t *testing.T) {
	linter := NewTemplateLabelPromQLRule()

	for _, tc := range []struct {
		name      string
		result    Result
		dashboard dashboard.Dashboard
	}{
		{
			name:   "Don't fail on non prometheus template.",
			result: ResultSuccess,
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Query: &dashboard.StringOrMap{String: toPtr("foo")},
						},
					},
				},
			},
		},
		{
			name:   "OK",
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
							Name:       "namespaces",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Query:      &dashboard.StringOrMap{String: toPtr("label_values(up{job=~\"$job\"}, namespace)")},
							Type:       "query",
							Label:      toPtr("Job"),
						},
					},
				},
			},
		},
		{
			name: "Error",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' template 'namespaces' invalid templated label 'label_values(up{, namespace)': 1:4: parse error: unexpected "," in label matching, expected identifier or "}"`,
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
							Name:       "namespaces",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Query:      &dashboard.StringOrMap{String: toPtr("label_values(up{, namespace)")},
							Type:       "query",
							Label:      toPtr("Job"),
						},
					},
				},
			},
		},
		{
			name: "Invalid function.",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' template 'namespaces' invalid templated label 'foo(up, namespace)': invalid 'function': foo`,
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
							Name:       "namespaces",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Query:      &dashboard.StringOrMap{String: toPtr("foo(up, namespace)")},
							Type:       "query",
							Label:      toPtr("Job"),
						},
					},
				},
			},
		},
		{
			name: "Invalid query expression.",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' template 'namespaces' invalid templated label 'foo': invalid 'query': foo`,
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
							Name:       "namespaces",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Query:      &dashboard.StringOrMap{String: toPtr("foo")},
							Type:       "query",
							Label:      toPtr("job"),
						},
					},
				},
			},
		},
		// Support main grafana variables.
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
							Name:       "namespaces",
							Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
							Query:      &dashboard.StringOrMap{String: toPtr("query_result(max by(namespaces) (max_over_time(memory{}[$__range])))")},
							Type:       "query",
							Label:      toPtr("job"),
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testRule(t, linter, tc.dashboard, tc.result)
		})
	}
}
