package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestTemplateDatasource(t *testing.T) {
	linter := NewTemplateDatasourceRule()

	for _, tc := range []struct {
		name      string
		result    []Result
		dashboard dashboard.Dashboard
	}{
		// 0 Data Sources
		{
			name: "0 Data Sources",
			result: []Result{{
				Severity: Error,
				Message:  "Dashboard 'test' does not have a templated data source",
			}},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
			},
		},
		// 1 Data Source
		{
			name: "1 Data Source",
			result: []Result{
				{
					Severity: Error,
					Message:  "Dashboard 'test' templated data source variable named 'foo', should be named '_datasource', or 'datasource'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled '', should be labeled ' data source', or 'Data source'",
				},
			},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type: "datasource",
							Name: "foo",
						},
					},
				},
			},
		},
		{
			name: "wrong name",
			result: []Result{
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'bar', should be labeled 'Bar data source', or 'Data source'",
				},
			},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Label: toPtr("bar"),
							Query: &dashboard.StringOrMap{String: toPtr("bar")},
						},
					},
				},
			},
		},
		{
			name:   "OK - Data source ",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
					},
				},
			},
		},
		{
			name:   "OK - Prometheus data source",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Label: toPtr("Prometheus data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
					},
				},
			},
		},
		{
			name:   "OK - name: prometheus_datasource",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "prometheus_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
					},
				},
			},
		},
		{
			name:   "OK - name: prometheus_datasource, label: Prometheus data source",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "prometheus_datasource",
							Label: toPtr("Prometheus data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
					},
				},
			},
		},
		{
			name:   "OK - name: loki_datasource, query: loki",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "loki_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("loki")},
						},
					},
				},
			},
		},
		{
			name:   "OK - name: datasource, query: loki",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("loki")},
						},
					},
				},
			},
		},
		// 2 or more Data Sources
		{
			name: "3 Data Sources - 0",
			result: []Result{
				{
					Severity: Error,
					Message:  "Dashboard 'test' templated data source variable named 'datasource', should be named 'prometheus_datasource'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Prometheus data source'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Loki data source'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Influx data source'",
				},
			},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Type:  "datasource",
							Name:  "loki_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("loki")},
						},
						{
							Type:  "datasource",
							Name:  "influx_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("influx")},
						},
					},
				},
			},
		},
		{
			name: "3 Data Sources - 1",
			result: []Result{
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Prometheus data source'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Loki data source'",
				},
				{
					Severity: Warning,
					Message:  "Dashboard 'test' templated data source variable labeled 'Data source', should be labeled 'Influx data source'",
				},
			},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "prometheus_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Type:  "datasource",
							Name:  "loki_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("loki")},
						},
						{
							Type:  "datasource",
							Name:  "influx_datasource",
							Label: toPtr("Data source"),
							Query: &dashboard.StringOrMap{String: toPtr("influx")},
						},
					},
				},
			},
		},
		{
			name:   "3 Data Sources - 2",
			result: []Result{ResultSuccess},
			dashboard: dashboard.Dashboard{
				Title: toPtr("test"),
				Templating: dashboard.DashboardDashboardTemplating{
					List: []dashboard.VariableModel{
						{
							Type:  "datasource",
							Name:  "prometheus_datasource",
							Label: toPtr("Prometheus data source"),
							Query: &dashboard.StringOrMap{String: toPtr("prometheus")},
						},
						{
							Type:  "datasource",
							Name:  "loki_datasource",
							Label: toPtr("Loki data source"),
							Query: &dashboard.StringOrMap{String: toPtr("loki")},
						},
						{
							Type:  "datasource",
							Name:  "influx_datasource",
							Label: toPtr("Influx data source"),
							Query: &dashboard.StringOrMap{String: toPtr("influx")},
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
