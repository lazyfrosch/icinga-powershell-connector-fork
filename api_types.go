package main

import "strings"

type APICheckResults map[string]APICheckResult

type APICheckResult struct {
	ExitCode    int
	CheckResult string
	Perfdata    []string
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
