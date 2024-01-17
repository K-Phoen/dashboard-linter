package lint

import (
	"encoding/json"

	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

type Severity int

const (
	Success Severity = iota
	Exclude
	Quiet
	Warning
	Error
	Fixed

	Prometheus = "prometheus"
	Loki       = "loki"
)

func NewDashboard(buf []byte) (dashboard.Dashboard, error) {
	var dash dashboard.Dashboard
	if err := json.Unmarshal(buf, &dash); err != nil {
		return dash, err
	}
	return dash, nil
}

// Target decorates cog's definition of a target to maintain the `idx` field,
// which is used to uniquely identify panel targets while linting.
type Target struct {
	Idx      int `json:"-"` // This is the only (best?) way to uniquely identify a target, it is set by GetPanels
	Original cogvariants.Dataquery
}
