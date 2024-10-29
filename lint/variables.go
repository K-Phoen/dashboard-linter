package lint

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

// https://grafana.com/docs/grafana/latest/variables/variable-types/global-variables/
var globalVariables = map[string]interface{}{
	"__rate_interval": "8869990787ms",
	"__interval":      "4867856611ms",
	"__interval_ms":   "7781188786",
	"__range_ms":      "6737667980",
	"__range_s":       "9397795485",
	"__range":         "6069770749ms",
	"__dashboard":     "AwREbnft",
	"__from":          time.Date(2020, 7, 13, 20, 19, 9, 254000000, time.UTC),
	"__to":            time.Date(2020, 7, 13, 20, 19, 9, 254000000, time.UTC),
	"__name":          "name",
	"__org":           42,
	"__org.name":      "orgname",
	"__user.id":       42,
	"__user.login":    "user",
	"__user.email":    "user@test.com",
	"timeFilter":      "time > now() - 7d",
	"__timeFilter":    "time > now() - 7d",
	"__auto":          "12345ms",
}

func getTemplateByType(d dashboard.Dashboard, t string) []dashboard.VariableModel {
	var retval []dashboard.VariableModel
	for _, templ := range d.Templating.List {
		if strings.EqualFold(string(templ.Type), t) {
			retval = append(retval, templ)
		}
	}
	return retval
}

func getTemplate(d dashboard.Dashboard, name string) *dashboard.VariableModel {
	for _, template := range d.Templating.List {
		if template.Name == name {
			return &template
		}
	}
	return nil
}

func getTemplateDatasource(d dashboard.Dashboard) *dashboard.VariableModel {
	for _, template := range d.Templating.List {
		if template.Type != "datasource" {
			continue
		}
		return &template
	}
	return nil
}

func stringValue(name string, value interface{}, kind, format string) (string, error) {
	switch val := value.(type) {
	case int:
		return strconv.Itoa(val), nil
	case time.Time:
		// Implements https://grafana.com/docs/grafana/latest/variables/variable-types/global-variables/#__from-and-__to
		switch kind {
		case "date":
			switch format {
			case "seconds":
				return strconv.FormatInt(val.Unix(), 10), nil
			case "iso":
				return val.Format(time.RFC3339), nil
			default:
				return "", fmt.Errorf("Unsupported momentjs time format: %s", format)
			}
		default:
			switch format {
			case "date":
				return val.Format(time.RFC3339), nil
			default:
				return strconv.FormatInt(val.UnixMilli(), 10), nil
			}
		}
	default:
		// Use variable name as sample value
		svalue := fmt.Sprintf("%s", value)
		// For list types, repeat it 3 times (arbitrary value)
		svalueList := []string{svalue, svalue, svalue}
		// Implements https://grafana.com/docs/grafana/latest/variables/advanced-variable-format-options/
		switch format {
		case "csv":
			return strings.Join(svalueList, ","), nil
		case "doublequote":
			return "\"" + strings.Join(svalueList, "\",\"") + "\"", nil
		case "glob":
			return "{" + strings.Join(svalueList, ",") + "}", nil
		case "json":
			data, err := json.Marshal(svalueList)
			if err != nil {
				return "", err
			}
			return string(data), nil
		case "lucene":
			return "(\"" + strings.Join(svalueList, "\" OR \"") + "\")", nil
		case "percentencode":
			return url.QueryEscape(strings.Join(svalueList, ",")), nil
		case "pipe":
			return strings.Join(svalueList, "|"), nil
		case "raw":
			return strings.Join(svalueList, ","), nil
		case "regex":
			return strings.Join(svalueList, "|"), nil
		case "singlequote":
			return "'" + strings.Join(svalueList, "','") + "'", nil
		case "sqlstring":
			return "'" + strings.Join(svalueList, "','") + "'", nil
		case "text":
			return strings.Join(svalueList, " + "), nil
		case "queryparam":
			values := url.Values{}
			for _, svalue := range svalueList {
				values.Add("var-"+name, svalue)
			}
			return values.Encode(), nil
		default:
			return svalue, nil
		}
	}
}

func removeVariableByName(name string, variables []dashboard.VariableModel) []dashboard.VariableModel {
	vars := make([]dashboard.VariableModel, 0, len(variables))
	for _, v := range variables {
		if v.Name == name {
			continue
		}
		vars = append(vars, v)
	}
	return vars
}

func variableSampleValue(s string, variables []dashboard.VariableModel) (string, error) {
	var name, kind, format string
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		// No format
		name = s
	case 2:
		// Could be __from:date, variable:csv, ...
		name = parts[0]
		format = parts[1]
	case 3:
		// Could be __from:date:iso, ...
		name = parts[0]
		kind = parts[1]
		format = parts[2]
	default:
		return "", fmt.Errorf("unknown variable format: %s", s)
	}
	// If it is part of the globals, return a string representation of a sample value
	if value, ok := globalVariables[name]; ok {
		return stringValue(name, value, kind, format)
	}
	// If it is an auto interval variable, replace with a sample value of 10s
	if strings.HasPrefix(name, "__auto_interval") {
		return "10s", nil
	}
	// If it is a template variable and we have a value, we use it
	for _, v := range variables {
		if v.Name != name {
			continue
		}

		// if it has a current value, use it
		if v.Current != nil && v.Current.Value.String != nil && *v.Current.Value.String != "" {
			// Recursively expand, without the current variable to avoid infinite recursion
			return expandVariables(*v.Current.Value.String, removeVariableByName(name, variables))
		}
		// If it has options, use the first option
		if len(v.Options) > 0 {
			// Recursively expand, without the current variable to avoid infinite recursion
			return expandVariables(*v.Options[0].Value.String, removeVariableByName(name, variables))
		}
	}
	// Assume variable type is a string
	return stringValue(name, name, kind, format)
}

var variableRegexp = regexp.MustCompile(
	strings.Join([]string{
		`\$([[:word:]]+)`,    // $var syntax
		`\$\{([^}]+)\}`,      // ${var} syntax
		`\[\[([^\[\]]+)\]\]`, // [[var]] syntax
	}, "|"),
)

func expandVariables(expr string, variables []dashboard.VariableModel) (string, error) {
	parts := strings.Split(expr, "\"")
	for i, part := range parts {
		if i%2 == 1 {
			// Inside a double quote string, just add it
			continue
		}

		// Accumulator to store the processed submatches
		var subparts []string
		// Cursor indicates where we are in the part being processed
		cursor := 0
		for _, v := range variableRegexp.FindAllStringSubmatchIndex(part, -1) {
			// Add all until match starts
			subparts = append(subparts, part[cursor:v[0]])
			// Iterate on all the subgroups and find the one that matched
			for j := 2; j < len(v); j += 2 {
				if v[j] < 0 {
					continue
				}
				// Replace the match with sample value
				val, err := variableSampleValue(part[v[j]:v[j+1]], variables)
				if err != nil {
					return "", err
				}
				subparts = append(subparts, val)
			}
			// Move the start cursor at the end of the current match
			cursor = v[1]
		}
		// Add rest of the string
		subparts = append(subparts, parts[i][cursor:])
		// Merge all back into the parts
		parts[i] = strings.Join(subparts, "")
	}
	return strings.Join(parts, "\""), nil
}

func expandLogQLVariables(expr string, variables []dashboard.VariableModel) (string, error) {
	lines := strings.Split(expr, "\n")
	for i, line := range lines {
		parts := strings.Split(line, "\"")
		for j, part := range parts {
			if j%2 == 1 {
				// Inside a double quote string, just add it
				continue
			}

			// Accumulator to store the processed submatches
			var subparts []string
			// Cursor indicates where we are in the part being processed
			cursor := 0
			for _, v := range variableRegexp.FindAllStringSubmatchIndex(part, -1) {
				// Add all until match starts
				subparts = append(subparts, part[cursor:v[0]])
				// Iterate on all the subgroups and find the one that matched
				for k := 2; k < len(v); k += 2 {
					if v[k] < 0 {
						continue
					}
					// Replace the match with sample value
					val, err := variableSampleValue(part[v[k]:v[k+1]], variables)
					if err != nil {
						return "", err
					}
					// If the variable is within square brackets, remove the '$' prefix
					if strings.HasPrefix(part[v[0]-1:v[0]], "[") && strings.HasSuffix(part[v[1]:v[1]+1], "]") {
						val = strings.TrimPrefix(val, "$")
					}
					subparts = append(subparts, val)
				}
				// Move the start cursor at the end of the current match
				cursor = v[1]
			}
			// Add rest of the string
			subparts = append(subparts, part[cursor:])
			// Merge all back into the parts
			parts[j] = strings.Join(subparts, "")
		}
		lines[i] = strings.Join(parts, "\"")
	}
	result := strings.Join(lines, "\n")
	return result, nil
}
