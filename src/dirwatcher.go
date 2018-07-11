package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"strconv"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	DirName   string `default:"C:\\temp" help:"File name being monitored"`
	DoRecurse string `default:"false"` help:"Whether to monitor top level of directory or recursively walk the dir and its subdirs"
}

type dirWatcherFile struct {
	DWFInfo   os.FileInfo
	DWFPath		string
}

const (
	integrationName    = "com.newrelic.dirwatcher"
	integrationVersion = "0.1.0"
)

var args argumentList

func getFileList() []dirWatcherFile {
	var dwflist []dirWatcherFile
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
			dwflist = append(dwflist, dirWatcherFile{DWFInfo: finfo, DWFPath: path})
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return dwflist
	} else {
		flist, err := ioutil.ReadDir(args.DirName)
		for _, finfo := range flist {
			filename := filepath.ToSlash(args.DirName+"/"+finfo.Name())
			dwflist = append(dwflist, dirWatcherFile{DWFInfo: finfo, DWFPath: filename})
		}
		if err != nil {
			log.Fatal(err)
		}
		return dwflist
	}
}

func populateInventory(dwflist []dirWatcherFile, inventory sdk.Inventory) error {
	for _, dwfile := range dwflist {
		inventory.SetItem(dwfile.DWFPath, "fileSize", dwfile.DWFInfo.Size())
		inventory.SetItem(dwfile.DWFPath, "mode", dwfile.DWFInfo.Mode().String())
		inventory.SetItem(dwfile.DWFPath, "modTime", dwfile.DWFInfo.ModTime().Unix())
		inventory.SetItem(dwfile.DWFPath, "isDir", strconv.FormatBool(dwfile.DWFInfo.IsDir()))
	}
	return nil
}

func populateMetrics(dwflist []dirWatcherFile, integration *sdk.Integration) error {
	for _, dwfile := range dwflist {
		ms := integration.NewMetricSet("DirWatcher")
		ms.SetMetric("filePath", dwfile.DWFPath, metric.ATTRIBUTE)
		ms.SetMetric("fileSize", dwfile.DWFInfo.Size(), metric.GAUGE)
		ms.SetMetric("mode", dwfile.DWFInfo.Mode().String(), metric.ATTRIBUTE)
		ms.SetMetric("modTime", dwfile.DWFInfo.ModTime().Unix(), metric.GAUGE)
		ms.SetMetric("isDir", strconv.FormatBool(dwfile.DWFInfo.IsDir()), metric.ATTRIBUTE)
	}
	return nil
}

func main() {
	log.SetupLogging(args.Verbose)

	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	filelist := getFileList()

	if args.All || args.Inventory {
		fatalIfErr(populateInventory(filelist, integration.Inventory))
	}

	if args.All || args.Metrics {
		fatalIfErr(populateMetrics(filelist, integration))
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
