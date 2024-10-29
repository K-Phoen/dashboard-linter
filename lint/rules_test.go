package lint_test

import (
	"os"
	"testing"

	"github.com/grafana/dashboard-linter/lint"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/stretchr/testify/assert"
)

func TestCustomRules(t *testing.T) {
	sampleDashboard, err := os.ReadFile("testdata/dashboard.json")
	assert.NoError(t, err)

	for _, tc := range []struct {
		desc string
		rule lint.Rule
	}{
		{
			desc: "Should allow addition of dashboard rule",
			rule: lint.NewDashboardRuleFunc(
				"test-dashboard-rule", "Test dashboard rule",
				func(dashboard.Dashboard) lint.DashboardRuleResults {
					return lint.DashboardRuleResults{Results: []lint.DashboardResult{{
						Result: lint.Result{Severity: lint.Error, Message: "Error found"},
					}}}
				},
			),
		},
		{
			desc: "Should allow addition of panel rule",
			rule: lint.NewPanelRuleFunc(
				"test-panel-rule", "Test panel rule",
				func(d dashboard.Dashboard, p dashboard.PanelOrRowPanel) lint.PanelRuleResults {
					return lint.PanelRuleResults{Results: []lint.PanelResult{{
						Result: lint.Result{Severity: lint.Error, Message: "Error found"},
					}}}
				},
			),
		},
		{
			desc: "Should allow addition of target rule",
			rule: lint.NewTargetRuleFunc(
				"test-target-rule", "Test target rule",
				func(dashboard.Dashboard, dashboard.PanelOrRowPanel, lint.Target) lint.TargetRuleResults {
					return lint.TargetRuleResults{Results: []lint.TargetResult{{
						Result: lint.Result{Severity: lint.Error, Message: "Error found"},
					}}}
				},
			),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			rules := lint.RuleSet{}
			rules.Add(tc.rule)

			d, err := lint.NewDashboard(sampleDashboard)
			assert.NoError(t, err)

			results, err := rules.Lint([]dashboard.Dashboard{d})
			assert.NoError(t, err)

			// Validate the error was added
			assert.GreaterOrEqual(t, len(results.ByRule()[tc.rule.Name()]), 1)
			r := results.ByRule()[tc.rule.Name()][0].Result
			assert.Equal(t, lint.Result{Severity: lint.Error, Message: "Error found"}, r.Results[0].Result)
		})
	}
}

func TestFixableRules(t *testing.T) {
	sampleDashboard, err := os.ReadFile("testdata/dashboard.json")
	assert.NoError(t, err)

	rule := lint.NewDashboardRuleFunc(
		"test-fixable-rule", "Test fixable rule",
		func(d dashboard.Dashboard) lint.DashboardRuleResults {
			rr := lint.DashboardRuleResults{}
			rr.AddFixableError(d, "fixing first issue", func(d *dashboard.Dashboard) {
				d.Title = cog.ToPtr(*d.Title + " fixed-once")
			})
			rr.AddFixableError(d, "fixing second issue", func(d *dashboard.Dashboard) {
				d.Title = cog.ToPtr(*d.Title + " fixed-twice")
			})
			return rr
		},
	)

	rules := lint.RuleSet{}
	rules.Add(rule)

	dash, err := lint.NewDashboard(sampleDashboard)
	assert.NoError(t, err)

	results, err := rules.Lint([]dashboard.Dashboard{dash})
	assert.NoError(t, err)

	results.AutoFix(&dash)

	assert.Equal(t, "Sample dashboard fixed-once fixed-twice", *dash.Title)
}
