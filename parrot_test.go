package parrot_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/oliveiraethales/parrot"
)

func TestParrot(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, parrot.Analyzer, "a")
}
