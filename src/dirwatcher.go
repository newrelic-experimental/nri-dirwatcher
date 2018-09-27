package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// ByteSize Constant for Megabytes
const (
	MB float64           = 1 << (10 * 2)
	G  metric.SourceType = metric.GAUGE
	A  metric.SourceType = metric.ATTRIBUTE
)

// File extends os.FileInfo
type File struct {
	os.FileInfo
	Path string
}

// ModInMinutes TODO
func (f *File) ModInMinutes() int64 {
	return int64(time.Now().UTC().Sub(f.ModTime()).Minutes())
}

// Result TODO description
type Result struct {
	ParentDir           string
	SizeMb              float64
	DirectoryCount      int
	FileCount           int
	LastModified        time.Time
	LastModifiedMinutes int64
	Oldest              File
	Newest              File
	LastError           error
	Errors              []string
	ErrorCount          int64
	Recursive           bool
}

// Debug prints it
func (r *Result) Debug() {
	fmt.Printf("\n%+v\n", r)
	fmt.Printf("Last Modified: %d\n", r.LastModified.Unix())
	fmt.Printf("Minutes since mod: %d\n", r.LastModifiedMinutes)
	fmt.Printf("Total SizeMb: %f\n", r.SizeMb)
	fmt.Printf("Dir Count: %d\n", r.DirectoryCount)
	fmt.Printf("File Count: %d\n", r.FileCount)
	fmt.Printf("Oldest File: %s, Minutes since mod: %d\n", r.Oldest.Name(), r.Oldest.ModInMinutes())
	fmt.Printf("Newest File: %s, Minutes since mod: %d\n", r.Newest.Name(), r.Newest.ModInMinutes())
}

// PopulateMetrics is method to impliment Infra SDK MetricSet
func (r *Result) PopulateMetrics(e *integration.Entity) (err error) {

	for _, list := range r.Errors {
		errorSet := e.NewMetricSet("DirWatcherErrors")
		errorSet.SetMetric("dir.target", r.ParentDir, A)
		errorSet.SetMetric("dir.error", list, A)
	}

	summarySet := e.NewMetricSet("DirWatcher")

	summarySet.SetMetric("dir.parent",
		r.ParentDir, A)

	summarySet.SetMetric("dir.recursive",
		strconv.FormatBool(r.Recursive), A)

	summarySet.SetMetric("dir.sizemb",
		r.SizeMb, G)

	summarySet.SetMetric("dir.directory_count",
		r.DirectoryCount, G)

	summarySet.SetMetric("dir.file_count",
		r.FileCount, G)

	summarySet.SetMetric("dir.last_modified",
		r.LastModified.Unix(), G)

	summarySet.SetMetric("dir.last_mod_minutes",
		r.LastModifiedMinutes, G)

	// Handle Empty Directories
	if r.Oldest.FileInfo == nil || r.Newest.FileInfo == nil {
		return
	}

	summarySet.SetMetric("dir.oldest.name",
		r.Oldest.Name(), A)

	summarySet.SetMetric("dir.oldest.last_modified",
		r.Oldest.ModTime().Unix(), G)

	summarySet.SetMetric("dir.oldest.last_mod_minutes",
		r.Oldest.ModInMinutes(), G)

	summarySet.SetMetric("dir.oldest.path",
		filepath.Dir(r.Oldest.Path), A)

	summarySet.SetMetric("dir.oldest.fullpath",
		r.Oldest.Path, A)

	summarySet.SetMetric("dir.oldest.sizebytes",
		r.Oldest.Size(), G)

	summarySet.SetMetric("dir.newest.name",
		r.Newest.Name(), A)

	summarySet.SetMetric("dir.newest.last_modified",
		r.Newest.ModTime().Unix(), G)

	summarySet.SetMetric("dir.newest.last_mod_minutes",
		r.Newest.ModInMinutes(), G)

	summarySet.SetMetric("dir.newest.path",
		filepath.Dir(r.Newest.Path), A)

	summarySet.SetMetric("dir.newest.fullpath",
		r.Newest.Path, A)

	summarySet.SetMetric("dir.newest.sizebytes",
		r.Newest.Size(), G)

	summarySet.SetMetric("dir.error_count",
		r.ErrorCount, G)

	if r.LastError != nil {
		summarySet.SetMetric("dir.last_error",
			r.LastError.Error(), A)
	}

	return
}

// Get TODO description
func Get(p string, recurse bool) *Result {

	// Absolute path required
	if f, err := os.Stat(p); err != nil || f.IsDir() == false {
		fmt.Fprintf(os.Stderr, "%s not an absolute path", p)
	}

	r := &Result{ParentDir: p, Errors: []string{}, Recursive: recurse}

	if targetDir, err := os.Stat(p); err == nil {
		d := File{targetDir, p}
		r.LastModified = d.ModTime()
		r.LastModifiedMinutes = d.ModInMinutes()
	}

	var err error

	if recurse {
		err = Walk(r)
	} else {
		err = ListContents(r)
	}

	if err != nil {
		// should never happen
		log.Fatalf("error walking directory: %v", err)
		os.Exit(1)
	}

	return r

}

// ListContents ...
func ListContents(r *Result) (err error) {
	r.DirectoryCount = 1
	files, err := ioutil.ReadDir(r.ParentDir)

	if err != nil {
		return err
	}

	if len(files) < 1 {
		return
	}

	for _, f := range files[1:] {

		if f.IsDir() {
			r.DirectoryCount++
			continue
		}

		if r.Newest.FileInfo == nil || r.Oldest.FileInfo == nil {
			r.Newest = File{f, r.ParentDir}
			r.Oldest = File{f, r.ParentDir}
			continue
		}

		if f.ModTime().After(r.Newest.ModTime()) {
			r.Newest = File{f, r.ParentDir}
		}

		if f.ModTime().Before(r.Oldest.ModTime()) {
			r.Oldest = File{f, r.ParentDir}
		}

		r.FileCount++
		r.SizeMb += float64(f.Size()) / MB

	}

	return

}

// Walk ...
func Walk(r *Result) error {

	return filepath.Walk(r.ParentDir, func(path string, f os.FileInfo, err error) error {

		if err != nil {
			r.LastError = err
			r.Errors = append(r.Errors, err.Error())
			r.ErrorCount++
			return nil
		}

		if f.IsDir() {
			r.DirectoryCount++
			return nil
		}

		// instantiate os.FileInfo interface on first file found,
		// otherwise conditions below panic on nil interface
		if r.Newest.FileInfo == nil || r.Oldest.FileInfo == nil {
			r.Newest = File{f, path}
			r.Oldest = File{f, path}
			return nil
		}

		if f.ModTime().After(r.Newest.ModTime()) {
			r.Newest = File{f, path}
		}

		if f.ModTime().Before(r.Oldest.ModTime()) {
			r.Oldest = File{f, path}
		}

		r.FileCount++
		r.SizeMb += float64(f.Size()) / MB

		return nil

	})

}
