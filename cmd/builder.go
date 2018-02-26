package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cryptopay-dev/yaga/cli"
)

const pathProjectExample = "github.com/cryptopay-dev/yaga/cmd/project_example"

var (
	projectName string
	projectPath string
)

func isEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}

func verifyOrCreateWorkdir(workdir string) error {
	mode, err := os.Stat(workdir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(workdir, 0755); err != nil {
			return err
		}
		return nil
	}

	if !mode.IsDir() {
		return fmt.Errorf("workdir should be directory")
	}

	isEmpty, err := isEmptyDir(workdir)
	if err != nil {
		return err
	}
	if !isEmpty {
		return fmt.Errorf("workdir should be empty")
	}

	return nil
}

func copyFiles(dst, src string) error {
	cpCmd := exec.Command("cp", "-rf", src, dst)
	return cpCmd.Run()
}

func projectBuilder(opts *cli.Options, dir string) (err error) {
	if len(dir) == 0 {
		return fmt.Errorf("workdir is undefined")
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}

	path := os.Getenv("GOPATH")
	if len(path) == 0 {
		return fmt.Errorf("GOPATH cannot be empty")
	}

	paths := strings.Split(path, ":")
	srcProjectPath, err := getPathProjectExample(paths)
	if err != nil {
		return err
	}

	validPath := false
	for _, path = range paths {
		path, _ = filepath.Abs(path + "/src")
		if strings.HasPrefix(dir, path) {
			validPath = true
			projectPath = strings.TrimLeft(dir, path)
			projectName = filepath.Base(projectPath)
			break
		}
	}
	if !validPath {
		return fmt.Errorf("workdir should be inside GOPATH")
	}
	fmt.Printf("Project example path: %s\n\n", pathProjectExample)
	fmt.Printf("Project mame: %s\n", projectName)
	fmt.Printf("Project full name: %s\n", projectPath)
	fmt.Printf("Project path: %s\n", dir)

	if err = verifyOrCreateWorkdir(dir); err != nil {
		return err
	}

	if err = copyFiles(dir, srcProjectPath+"/"); err != nil {
		return err
	}

	parseGolangFiles(opts, dir, pathProjectExample, projectPath)

	fmt.Println("Done.")

	return nil
}

func isDirAndExist(dir string) (bool, error) {
	mode, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return mode.IsDir(), nil
}

func getPathProjectExample(paths []string) (string, error) {
	for _, path := range paths {
		path, _ = filepath.Abs(filepath.Join(path, "src", pathProjectExample))
		if ok, err := isDirAndExist(path); err != nil {
			return "", err
		} else if ok {
			return path, nil
		}
	}

	return "", fmt.Errorf("cannot found project example")
}
