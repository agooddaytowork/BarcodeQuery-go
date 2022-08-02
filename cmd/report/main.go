package main

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
)

func main() {
	beam.Init()

	//// Create the Pipeline object and root scope.
	//pipeline, scope := beam.NewPipelineWithRoot()
	//
	//lines := textio.Read(scope, "gs://some/inputData.txt")
}
