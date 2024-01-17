package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestPanelTitleDescription(t *testing.T) {
	linter := NewPanelTitleDescriptionRule()

	for _, tc := range []struct {
		result []Result
		panel  dashboard.Panel
	}{
		{
			result: []Result{
				{Severity: Error, Message: "Dashboard 'test', panel with id '1' has missing title"},
				{Severity: Error, Message: "Dashboard 'test', panel with id '1' has missing description"},
			},
			panel: dashboard.Panel{
				Type:        "singlestat",
				Id:          toPtr(uint32(1)),
				Title:       toPtr(""),
				Description: toPtr(""),
			},
		},
		{
			result: []Result{
				{Severity: Error, Message: "Dashboard 'test', panel 'title' has missing description"},
			},
			panel: dashboard.Panel{
				Type:        "singlestat",
				Id:          toPtr(uint32(2)),
				Title:       toPtr("title"),
				Description: toPtr(""),
			},
		},
		{
			result: []Result{
				{Severity: Error, Message: "Dashboard 'test', panel with id '3' has missing title"},
			},
			panel: dashboard.Panel{
				Type:        "singlestat",
				Id:          toPtr(uint32(3)),
				Title:       toPtr(""),
				Description: toPtr("description"),
			},
		},
		{
			result: []Result{ResultSuccess},
			panel: dashboard.Panel{
				Type:        "singlestat",
				Id:          toPtr(uint32(4)),
				Title:       toPtr("testpanel"),
				Description: toPtr("testdescription"),
			},
		},
	} {
		panels := []dashboard.PanelOrRowPanel{
			{Panel: &tc.panel},
		}
		testMultiResultRule(t, linter, dashboard.Dashboard{Title: toPtr("test"), Panels: panels}, tc.result)
	}
}
