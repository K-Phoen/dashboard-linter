package lint

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func TestPanelUnits(t *testing.T) {
	linter := NewPanelUnitsRule()

	for _, tc := range []struct {
		name   string
		result Result
		panel  dashboard.Panel
	}{
		{
			name: "invalid unit",
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test', panel 'bar' has no or invalid units defined: 'MyInvalidUnit'",
			},
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
				FieldConfig: &dashboard.FieldConfigSource{
					Defaults: dashboard.FieldConfig{
						Unit: toPtr("MyInvalidUnit"),
					},
				},
			},
		},
		{
			name: "missing FieldConfig",
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test', panel 'bar' has no or invalid units defined: ''",
			},
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
			},
		},
		{
			name: "empty FieldConfig",
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test', panel 'bar' has no or invalid units defined: ''",
			},
			panel: dashboard.Panel{
				Type:        "singlestat",
				Datasource:  &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:       toPtr("bar"),
				FieldConfig: &dashboard.FieldConfigSource{},
			},
		},
		{
			name:   "valid",
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
				FieldConfig: &dashboard.FieldConfigSource{
					Defaults: dashboard.FieldConfig{
						Unit: toPtr("short"),
					},
				},
			},
		},
		{
			name:   "none - scalar",
			result: ResultSuccess,
			panel: dashboard.Panel{
				Type:       "singlestat",
				Datasource: &dashboard.DataSourceRef{Uid: toPtr("foo")},
				Title:      toPtr("bar"),
				FieldConfig: &dashboard.FieldConfigSource{
					Defaults: dashboard.FieldConfig{
						Unit: toPtr("none"),
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			panels := []dashboard.PanelOrRowPanel{
				{Panel: &tc.panel},
			}
			testRule(t, linter, dashboard.Dashboard{Title: toPtr("test"), Panels: panels}, tc.result)
		})
	}
}
