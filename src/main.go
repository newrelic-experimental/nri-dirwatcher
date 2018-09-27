package main

import (
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Dirwatch        string `help:"Directory to watch"`
	DirwatchRecurse bool   `default:"true" help:"Recurse dirwatch path"`
}

const (
	integrationName    = "com.newrelic.mountstats"
	integrationVersion = "0.1.0"
)

// Args is args ...
var Args argumentList

func main() {

	i, err := integration.New(integrationName, integrationVersion, integration.Args(&Args))
	if err != nil {
		panic(err)
	}

	e := i.LocalEntity()

	if len(Args.Dirwatch) > 0 {
		for _, path := range strings.Split(Args.Dirwatch, ",") {
			Get(strings.TrimSpace(path), Args.DirwatchRecurse).
				PopulateMetrics(e)
		}
	}

	if ok := i.Publish(); ok != nil {
		panic(ok)
	}

}
