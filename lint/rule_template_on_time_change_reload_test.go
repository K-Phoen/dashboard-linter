package lint

import (
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/stretchr/testify/require"
)

func TestTemplateOnTimeRangeReloadRule(t *testing.T) {
	linter := NewTemplateOnTimeRangeReloadRule()

	good := []dashboard.VariableModel{
		{
			Type:  "datasource",
			Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
		},
		{
			Name:       "namespaces",
			Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
			Query:      &dashboard.StringOrMap{String: toPtr("label_values(up{job=~\"$job\"}, namespace)")},
			Type:       "query",
			Label:      toPtr("job"),
			Refresh:    toPtr(dashboard.VariableRefreshOnTimeRangeChanged),
		},
	}
	for _, tc := range []struct {
		name      string
		result    Result
		dashboard dashboard.Dashboard
		fixed     *dashboard.Dashboard
	}{
		{
			name:   "OK",
			result: ResultSuccess,
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: good,
				},
			},
		},
		{
			name: "autofix",
			result: Result{
				Severity: Fixed,
				Message:  `Dashboard 'test' templated datasource variable named 'namespaces', should be set to be refreshed 'On Time Range Change (value 2)', is currently '1'`,
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
							Query:      &dashboard.StringOrMap{String: toPtr("label_values(up{job=~\"$job\"}, namespace)")},
							Type:       "query",
							Label:      toPtr("job"),
							Refresh:    toPtr(dashboard.VariableRefreshOnDashboardLoad),
						},
					},
				},
			},
			fixed: &dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: good,
				},
			},
		},
		{
			name: "error",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' templated datasource variable named 'namespaces', should be set to be refreshed 'On Time Range Change (value 2)', is currently '1'`,
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
							Query:      &dashboard.StringOrMap{String: toPtr("label_values(up{job=~\"$job\"}, namespace)")},
							Type:       "query",
							Label:      toPtr("job"),
							Refresh:    toPtr(dashboard.VariableRefreshOnDashboardLoad),
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			autofix := tc.fixed != nil
			testRuleWithAutofix(t, linter, &tc.dashboard, []Result{tc.result}, autofix)
			if autofix {
				expected, _ := json.Marshal(tc.fixed)
				actual, _ := json.Marshal(tc.dashboard)
				require.Equal(t, string(expected), string(actual))
			}
		})
	}
}
