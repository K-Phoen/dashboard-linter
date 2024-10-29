package lint

import (
	"testing"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func TestPanelNoTargets(t *testing.T) {
	linter := NewPanelNoTargetsRule()

	for _, tc := range []struct {
		result Result
		panel  dashboard.Panel
	}{
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test', panel 'bar' has no targets",
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
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
				Targets: []cogvariants.Dataquery{
					prometheus.Dataquery{
						Expr: `sum(rate(foo[5m]))`,
					},
				},
			},
		},
	} {
		panels := []dashboard.PanelOrRowPanel{
			{Panel: &tc.panel},
		}
		testRule(t, linter, dashboard.Dashboard{Title: toPtr("test"), Panels: panels}, tc.result)
	}
}
