package lint

import (
	"fmt"
	"os"
	"sort"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

var ResultSuccess = Result{
	Severity: Success,
	Message:  "OK",
}

type Result struct {
	Severity Severity
	Message  string
}

type FixableResult struct {
	Result
	Fix func(*dashboard.Dashboard) // if nil, it cannot be fixed
}

type RuleResults struct {
	Results []FixableResult
}

type TargetResult struct {
	Result
	Fix func(dashboard.Dashboard, dashboard.PanelOrRowPanel, Target)
}

type TargetRuleResults struct {
	Results []TargetResult
}

func (r *TargetRuleResults) AddError(d dashboard.Dashboard, p dashboard.PanelOrRowPanel, t Target, message string) {
	title := ""
	if p.RowPanel != nil && p.RowPanel.Title != nil {
		title = *p.RowPanel.Title
	} else if p.Panel != nil && p.Panel.Title != nil {
		title = *p.Panel.Title
	}
	r.Results = append(r.Results, TargetResult{
		Result: Result{
			Severity: Error,
			Message:  fmt.Sprintf("Dashboard '%s', panel '%s', target idx '%d' %s", *d.Title, title, t.Idx, message),
		},
	})
}

type PanelResult struct {
	Result
	Fix func(dashboard.Dashboard, *dashboard.PanelOrRowPanel)
}

type PanelRuleResults struct {
	Results []PanelResult
}

func (r *PanelRuleResults) AddError(d dashboard.Dashboard, p dashboard.PanelOrRowPanel, message string) {
	var id uint32
	title := ""
	if p.RowPanel != nil {
		id = p.RowPanel.Id
		if p.RowPanel.Title != nil {
			title = *p.RowPanel.Title
		}
	} else if p.Panel != nil {
		if p.Panel.Id != nil {
			id = *p.Panel.Id
		}
		if p.Panel.Title != nil {
			title = *p.Panel.Title
		}
	}

	msg := fmt.Sprintf("Dashboard '%s', panel '%s' %s", *d.Title, title, message)
	if title == "" {
		msg = fmt.Sprintf("Dashboard '%s', panel with id '%d' %s", *d.Title, id, message)
	}

	r.Results = append(r.Results, PanelResult{
		Result: Result{
			Severity: Error,
			Message:  msg,
		},
	})
}

type DashboardResult struct {
	Result
	Fix func(*dashboard.Dashboard)
}

type DashboardRuleResults struct {
	Results []DashboardResult
}

func dashboardMessage(d dashboard.Dashboard, message string) string {
	return fmt.Sprintf("Dashboard '%s' %s", *d.Title, message)
}

func (r *DashboardRuleResults) AddError(d dashboard.Dashboard, message string) {
	r.Results = append(r.Results, DashboardResult{
		Result: Result{
			Severity: Error,
			Message:  dashboardMessage(d, message),
		},
	})
}

func (r *DashboardRuleResults) AddFixableError(d dashboard.Dashboard, message string, fix func(*dashboard.Dashboard)) {
	r.Results = append(r.Results, DashboardResult{
		Result: Result{
			Severity: Error,
			Message:  dashboardMessage(d, message),
		},
		Fix: fix,
	})
}

func (r *DashboardRuleResults) AddWarning(d dashboard.Dashboard, message string) {
	r.Results = append(r.Results, DashboardResult{
		Result: Result{
			Severity: Warning,
			Message:  dashboardMessage(d, message),
		},
	})
}

// ResultContext is used by ResultSet to keep all the state data about a lint execution and it's results.
type ResultContext struct {
	Result    RuleResults
	Rule      Rule
	Dashboard *dashboard.Dashboard
	Panel     *dashboard.PanelOrRowPanel
	Target    *Target
}

func (r Result) TtyPrint() {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	var Orange = "\033[38;5;208m"
	var sym string
	switch s := r.Severity; s {
	case Success:
		sym = Green + "✔️" + Reset
	case Fixed:
		sym = Orange + "🛠️ (fixed)" + Reset
	case Exclude:
		sym = "➖"
	case Warning:
		sym = Yellow + "⚠️" + Reset
	case Error:
		sym = Red + "❌" + Reset
	case Quiet:
		return
	}

	fmt.Fprintf(os.Stdout, "[%s] %s\n", sym, r.Message)
}

type ResultSet struct {
	results []ResultContext
	config  *ConfigurationFile
}

// Configure adds, and applies the provided configuration to all results currently in the ResultSet
func (rs *ResultSet) Configure(c *ConfigurationFile) {
	rs.config = c
	for i := range rs.results {
		rs.results[i] = rs.config.Apply(rs.results[i])
	}
}

// AddResult adds a result to the ResultSet, applying the current configuration if set
func (rs *ResultSet) AddResult(r ResultContext) {
	if rs.config != nil {
		r = rs.config.Apply(r)
	}
	rs.results = append(rs.results, r)
}

func (rs *ResultSet) MaximumSeverity() Severity {
	retVal := Success
	for _, res := range rs.results {
		for _, r := range res.Result.Results {
			if r.Severity > retVal {
				retVal = r.Severity
			}
		}
	}
	return retVal
}

func (rs *ResultSet) ByRule() map[string][]ResultContext {
	ret := make(map[string][]ResultContext)
	for _, res := range rs.results {
		ret[res.Rule.Name()] = append(ret[res.Rule.Name()], res)
	}
	for _, rule := range ret {
		sort.SliceStable(rule, func(i, j int) bool {
			if rule[i].Dashboard.Title == nil || rule[j].Dashboard.Title == nil {
				return false
			}

			return *rule[i].Dashboard.Title < *rule[j].Dashboard.Title
		})
	}
	return ret
}

func (rs *ResultSet) ReportByRule() {
	byRule := rs.ByRule()
	rules := make([]string, 0, len(byRule))
	for r := range byRule {
		rules = append(rules, r)
	}
	sort.Strings(rules)

	for _, rule := range rules {
		fmt.Fprintln(os.Stdout, byRule[rule][0].Rule.Description())
		for _, rr := range byRule[rule] {
			for _, r := range rr.Result.Results {
				if r.Severity == Exclude && !rs.config.Verbose {
					continue
				}
				r.TtyPrint()
			}
		}
	}
}

func (rs *ResultSet) AutoFix(d *dashboard.Dashboard) int {
	changes := 0
	for _, r := range rs.results {
		for i, fixableResult := range r.Result.Results {
			if fixableResult.Fix != nil {
				// Fix is only present when something can be fixed
				fixableResult.Fix(d)
				changes++
				r.Result.Results[i].Result.Severity = Fixed
			}
		}
	}
	return changes
}
