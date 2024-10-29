package lint

import (
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/stretchr/testify/require"
)

func TestNewUneditableRule(t *testing.T) {
	linter := NewUneditableRule()

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
				Title:    toPtr("test"),
				Editable: toPtr(false),
			},
		},
		{
			name: "error",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' is editable, it should be set to 'editable: false'`,
			},
			dashboard: dashboard.Dashboard{
				Title:    toPtr("test"),
				Editable: toPtr(true),
			},
		},
		{
			name: "error",
			result: Result{
				Severity: Error,
				Message:  `Dashboard 'test' is editable, it should be set to 'editable: false'`,
			},
			dashboard: dashboard.Dashboard{
				Title:    toPtr("test"),
				Editable: nil,
			},
		},
		{
			name: "autofix",
			result: Result{
				Severity: Fixed,
				Message:  `Dashboard 'test' is editable, it should be set to 'editable: false'`,
			},
			dashboard: dashboard.Dashboard{
				Title:    toPtr("test"),
				Editable: toPtr(true),
			},
			fixed: &dashboard.Dashboard{
				Title:    toPtr("test"),
				Editable: toPtr(false),
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
