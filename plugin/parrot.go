//go:build plugin

// Plugin for golangci-lint.
// Build with: go build -buildmode=plugin -o parrot.so ./plugin
package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/oliveiraethales/parrot"
)

var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{parrot.Analyzer}
}
