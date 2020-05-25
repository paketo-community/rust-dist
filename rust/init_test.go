package rust_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitRust(t *testing.T) {
	suite := spec.New("rust", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("LogEmitter", testLogEmitter)
	suite.Run(t)
}
