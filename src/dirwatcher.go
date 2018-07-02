package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	DirName   string `default:"C:\\temp" help:"File name being monitored"`
	DoRecurse string `default:"false"`
}

const (
	integrationName    = "com.newrelic.dirwatcher"
	integrationVersion = "0.1.0"
)

var args argumentList

func populateInventory(inventory sdk.Inventory) error {
	if strings.ToLower(args.DoRecurse) == "true" {
		fixedName := args.DirName
		if !strings.HasSuffix("/", args.DirName) {
			fixedName += "/"
		}
		err := filepath.Walk(fixedName, func(path string, finfo os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
				return nil
			}
			insertFileInfo(filepath.ToSlash(path), finfo, &inventory)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err := ioutil.ReadDir(args.DirName)
		if err != nil {
			log.Fatal(err)
		}
		for _, finfo := range files {
			insertFileInfo(filepath.ToSlash(args.DirName+"/"+finfo.Name()), finfo, &inventory)
		}
	}

	return nil
}

func insertFileInfo(filename string, fileinfo os.FileInfo, inventory *sdk.Inventory) {
	inventory.SetItem(filename, "fileSize", fileinfo.Size())
	inventory.SetItem(filename, "mode", fileinfo.Mode().String())
	inventory.SetItem(filename, "modTime", fileinfo.ModTime())
	inventory.SetItem(filename, "isDir", fileinfo.IsDir())
}

func populateMetrics(ms *metric.MetricSet) error {
	var fileSize int64 = 0
	fileCount := 0
	fixedName := args.DirName
	if !strings.HasSuffix("/", args.DirName) {
		fixedName += "/"
	}

	if strings.ToLower(args.DoRecurse) == "true" {
		err := filepath.Walk(fixedName, func(path string, finfo os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
				return nil
			}
			fileSize += finfo.Size()
			fileCount++
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err := ioutil.ReadDir(args.DirName)
		if err != nil {
			log.Fatal(err)
		}
		for _, finfo := range files {
			fileSize += finfo.Size()
			fileCount++
		}
	}
	ms.SetMetric("dirPath", fixedName, metric.ATTRIBUTE)
	ms.SetMetric("fileSize", fileSize, metric.GAUGE)
	ms.SetMetric("fileCount", fileCount, metric.GAUGE)
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
