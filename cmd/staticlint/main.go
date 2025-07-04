package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"strings"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"

	"honnef.co/go/tools/staticcheck"

	"github.com/SiyovushAbdulloev/metriks_sprint_1/cmd/staticlint/analyzers/noosexit"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		inspect.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
	)

	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, a.Analyzer)
			break
		}
	}

	analyzers = append(analyzers, noosexit.Analyzer)

	multichecker.Main(analyzers...)
}
