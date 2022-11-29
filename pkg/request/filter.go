package request

import (
	"net/url"
	"strconv"
	"strings"
)

var (
	comparisonOperators = map[string]string{
		"eq":   "=",
		"lte":  "<=",
		"gte":  ">=",
		"lt":   "<",
		"gt":   ">",
		"asc":  "asc",
		"desc": "desc",
	}
)

type FilteredResult struct {
	Field    string
	Type     string
	Value    interface{}
	Operator string
}

func QueryStringParser(queryStr string, filters map[string]string) []FilteredResult {
	type Map struct {
		Key   string
		Value string
	}

	params := make([]Map, 0)
	searchFilters := make([]FilteredResult, 0)

	parts := strings.Split(queryStr, "&")

	for _, part := range parts {
		split := strings.Split(part, "=")

		if len(split) > 1 && split[1] != "" {
			params = append(params, Map{
				Key:   split[0],
				Value: split[1],
			})
		} else {
			params = append(params, Map{
				Key:   split[0],
				Value: "",
			})
		}
	}

	for _, param := range params {
		for name, varType := range filters {
			if param.Key == name {
				esc, _ := url.QueryUnescape(param.Value)
				parseValue, operator := RHSParser(esc, varType)

				searchFilters = append(searchFilters, FilteredResult{
					Field:    param.Key,
					Type:     varType,
					Value:    parseValue,
					Operator: operator,
				})
				break
			}
		}
	}

	return searchFilters
}

func RHSParser(queryStrValue string, valueType string) (value interface{}, comparisonOperator string) {
	var val interface{}
	var cOperator = "="
	parts := strings.Split(queryStrValue, ":")

	l := len(parts)
	if valueType == "int" {
		var number int64
		number, _ = strconv.ParseInt(parts[0], 10, 64)
		val = number
	} else if valueType == "float" {
		number := 0.0
		number, _ = strconv.ParseFloat(parts[0], 64)
		val = number
	} else {
		val = parts[0]
	}

	if l == 1 {
		cOperator = comparisonOperators["eq"]
		return val, cOperator
	}

	if comparisonOperators[parts[1]] != "" {
		cOperator = comparisonOperators[parts[1]]
	}

	return val, cOperator
}
