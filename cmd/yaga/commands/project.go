package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var examplePath = "github.com/cryptopay-dev/yaga/cmd/yaga/project_example"

var postInfo = `Use it:

cd %s

// Prepare config:
cp config.example.yaml config.yaml
vi config.yaml

// Get dependencies:
dep ensure

// Run application:
go run main.go
`

func isEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}

func copyGoFile(apprelpath, src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()

	gout, err := ioutil.TempFile("", "go-")
	if err != nil {
		return
	}

	defer gout.Close()

	data, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}

	data = bytes.Replace(data, []byte(examplePath), []byte(apprelpath), -1)
	dataLen := len(data)

	var cnt int

	if cnt, err = gout.Write(data); err != nil {
		return
	} else if cnt != dataLen {
		err = fmt.Errorf("data not fully copied: %d != %d", cnt, dataLen)
		return
	}

	return copyFileContents(gout.Name(), dst)
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func newProject() cli.Command {
	action := func(ctx *cli.Context) error {
		var args = ctx.Args()

		if len(args) != 1 {
			errors("project path not set")
		}

		appath := args[0]
		gopath := os.Getenv("GOPATH")

		apprelpath := strings.Replace(appath, path.Join(gopath, "src"), "", -1)
		apprelpath = strings.Replace(apprelpath, "/", "", 1)
		apprelpath = strings.Replace(apprelpath, "\\", "", 1)
		apprelpath = strings.Replace(apprelpath, "\\", "/", -1) // fix to normal-path

		expath := path.Join(gopath, "src", examplePath)

		if _, err := os.Stat(expath); err != nil {
			errorsf("project template not found: %v", err)
		}

		info("Try to check project path")

		if !strings.Contains(appath, gopath) {
			errorsf("project path must be in GOPATH(%s)", gopath)
		}

		if file, err := os.Stat(appath); err != nil {
			if err = os.Mkdir(appath, 0700); err != nil {
				errorsf("Can't create project path: %v", err)
			}

			info("Project path created")
		} else if !file.IsDir() {
			errors("Project path not folder")
		} else {
			info("Project path found")
		}

		if ok, err := isEmptyDir(appath); err != nil {
			errorsf("Can't check project path: %v", err)
		} else if !ok {
			errors("Project path must be empty")
		}

		if err := filepath.Walk(expath, func(p string, i os.FileInfo, err error) error {
			if p == expath {
				return nil
			}

			relpath := strings.Replace(p, expath, "", -1)

			pp := path.Join(appath, relpath)
			if i.IsDir() {
				if _, err := os.Stat(pp); err != nil {
					if err = os.Mkdir(pp, 0700); err != nil {
						errorsf("Can't create path: %v", err)
					}

					infof("Create dir  %s", path.Join(apprelpath, relpath))
				}

				return nil
			}

			if path.Ext(p) != ".go" {
				if err := copyFileContents(p, pp); err != nil {
					errorsf("Can't copy file: %v", err)
				}
			} else if err := copyGoFile(apprelpath, p, pp); err != nil {
				errorsf("Can't copy file(%s): %v", i.Name(), err)
			}

			infof("Create file %s", path.Join(apprelpath, relpath))

			return nil
		}); err != nil {
			errorsf("Can't copy example project: %v", err)
		}

		print("Project created")
		printf(postInfo, appath)

		return nil
	}

	return cli.Command{
		Name:        "new",
		ShortName:   "n",
		Usage:       "new <work-dir>",
		Description: "Create new project",
		Before:      nil,
		After:       nil,
		Action:      action,
	}
}
