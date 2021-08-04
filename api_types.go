package main

import (
	"bytes"
	"encoding/json"
	"strings"
)

type APICheckResults map[string]APICheckResult

type APIPerfdataList []string

type APICheckResult struct {
	ExitCode    int
	CheckResult string
	Perfdata    APIPerfdataList
}

func (r APICheckResult) String() (s string) {
	s = strings.TrimSpace(r.CheckResult)

	if len(r.Perfdata) > 0 {
		s += "\n|"

		for _, p := range r.Perfdata {
			s += " " + strings.TrimSpace(p)
		}
	}

	s += "\n"

	return
}

// UnmarshalJSON makes sure we can de-serialize JSON.
//
// The API can return {} when no perfdata is set.
func (p *APIPerfdataList) UnmarshalJSON(data []byte) error {
	var value []string

	// catch empty object and return empty []string
	if bytes.Compare(data, []byte("{}")) == 0 {
		value = []string{}
	} else {
		err := json.Unmarshal(data, &value)
		if err != nil {
			return err
		}
	}

	*p = value

	return nil
}
