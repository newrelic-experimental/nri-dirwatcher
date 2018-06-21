package main

import (
	"io/ioutil"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
	// "fmt"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	// Filelocation string `default:"/Users/ayork/go/src/posready" help:"File location to monitor"`
	// Filename     string `default:"PosReady.flg" help:"File name being monitored"`
	DirName string `default:"C:\\Program Files\\New Relic\\newrelic-infra\\custom-integrations\\newrelic-infra-mssql" help:"File name being monitored"`
}

const (
	integrationName    = "com.newrelic.dirwatcher"
	integrationVersion = "0.1.0"
)

var args argumentList

func populateInventory(inventory sdk.Inventory) error {
	// Insert here the logic of your integration to get the inventory data
	// Ex: inventory.SetItem("softwareVersion", "value", "1.0.1")
	// --

	// files, err := ioutil.ReadDir(args.Filelocation)
	files, err := ioutil.ReadDir(args.DirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		inventory.SetItem(f.Name(), "value", "enabled")
	}

	return nil
}

func populateMetrics(ms *metric.MetricSet) error {
	// Insert here the logic of your integration to get the metrics data
	// files, err := ioutil.ReadDir(args.DirName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for _, f := range files {
	// 	// log.Debug(args.Filename + " Found")
	// 	// log.Debug("File Location: " + args.Filelocation)
	// 	// fmt.Println(f.Name())
	// 	// ms.SetMetric(f.Name(), "ENABLED", metric.ATTRIBUTE)
	// 	ms.SetMetric(f.Name(), "ENABLED", metric.ATTRIBUTE)
	// }

	return nil
}

func main() {
	log.SetupLogging(args.Verbose)

	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	if args.All || args.Inventory {
		fatalIfErr(populateInventory(integration.Inventory))
	}

	if args.All || args.Metrics {
		// fatalIfNotDefined(args.fileName, "Missing fileName parameter")
		ms := integration.NewMetricSet("DirWatcher")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
