package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/stretchr/testify/require"
)

func TestPanelDatasource(t *testing.T) {
	linter := NewPanelDatasourceRule()

	for _, tc := range []struct {
		result    Result
		panel     dashboard.Panel
		templates []dashboard.VariableModel
	}{
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test', panel 'bar' does not use a templated datasource, uses 'foo'",
			},
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
			},
		},
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
			},
			templates: []dashboard.VariableModel{
				{
					Type: "datasource",
					Name: "datasource",
				},
			},
		},
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("$datasource")},
			},
			templates: []dashboard.VariableModel{
				{
					Type: "datasource",
					Name: "datasource",
				},
			},
		},
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("$prometheus_datasource")},
			},
			templates: []dashboard.VariableModel{
				{
					Type: "datasource",
					Name: "prometheus_datasource",
				},
			},
		},
		{
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("${prometheus_datasource}")},
			},
			templates: []dashboard.VariableModel{
				{
					Type: "datasource",
					Name: "prometheus_datasource",
				},
			},
		},
	} {
		testRule(t, linter, dashboard.Dashboard{
			Title: toPtr("test"),
			Panels: []dashboard.PanelOrRowPanel{
				{Panel: &tc.panel},
			},
			Templating: dashboard.DashboardDashboardTemplating{List: tc.templates},
		}, tc.result)
	}
}

func toPtr[T any](input T) *T {
	return &input
}

// testRule is a small helper that tests a lint rule and expects it to only return
// a single result.
func testRule(t *testing.T, rule Rule, d dashboard.Dashboard, result Result) {
	testRuleWithAutofix(t, rule, &d, []Result{result}, false)
}
func testMultiResultRule(t *testing.T, rule Rule, d dashboard.Dashboard, result []Result) {
	testRuleWithAutofix(t, rule, &d, result, false)
}

func testRuleWithAutofix(t *testing.T, rule Rule, d *dashboard.Dashboard, result []Result, autofix bool) {
	rs := ResultSet{}
	rule.Lint(*d, &rs)
	if autofix {
		rs.AutoFix(d)
	}
	require.Len(t, rs.results, 1)
	actual := rs.results[0].Result
	if actual.Results[0].Severity == Quiet {
		// all test cases expect success
		actual.Results[0].Severity = Success
	}
	rr := make([]Result, len(actual.Results))
	for i, r := range actual.Results {
		rr[i] = r.Result
	}

	require.Equal(t, result, rr)
}
