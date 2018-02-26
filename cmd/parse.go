package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cryptopay-dev/yaga/cli"
)

func parseGolangFiles(opts *cli.Options, dirname, oldProjectName, newProjectName string) {
	waiter := &sync.WaitGroup{}
	filepath.Walk(dirname, makeWalkFunc(opts, oldProjectName, newProjectName, waiter))
	waiter.Wait()
}

func makeWalkFunc(opts *cli.Options, oldProjectName, newProjectName string, waiter *sync.WaitGroup) func(string, os.FileInfo, error) error {
	newProject := []byte(newProjectName)
	oldProject := []byte(oldProjectName)

	return func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Size() > 0 && filepath.Ext(info.Name()) == ".go" && info.Mode().IsRegular() {
			waiter.Add(1)
			go processFile(opts, path, info, oldProject, newProject, waiter)
		} else if err != nil {
			opts.Logger.Error(err)
		}

		return nil // we ignore all errors
	}
}

func processFile(opts *cli.Options, filename string, info os.FileInfo, oldProject, newProject []byte, waiter *sync.WaitGroup) {
	defer waiter.Done()

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		opts.Logger.Error(err)
	}

	buf = bytes.Replace(buf, oldProject, newProject, -1)

	err = ioutil.WriteFile(filename, buf, info.Mode())
	if err != nil {
		opts.Logger.Error(err)
	}
}
