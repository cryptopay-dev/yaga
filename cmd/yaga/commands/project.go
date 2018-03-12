package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/cryptopay-dev/yaga/logger"
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

func checkGopath(appath, gopath string) (bool, string) {
	var (
		err     error
		relpath string
		paths   = strings.Split(gopath, ":")
	)

	for _, p := range paths {
		p, err = filepath.Abs(p)
		if err != nil {
			continue
		}

		if strings.HasPrefix(appath, p) {
			relpath = strings.Replace(appath, p, "", -1)
			relpath = strings.Replace(relpath, "/", "", 1)    // for unix
			relpath = strings.Replace(relpath, "\\", "", 1)   // for windows
			relpath = strings.Replace(relpath, "\\", "/", -1) // fix to normal-path
			return true, relpath
		}
	}
	return false, relpath
}

func projectFormatter(log logger.Logger, appath string) {
	var tools = []string{"goimports", "gofmt"}
	for _, tool := range tools {
		var message = "OK"
		if err := exec.Command(tool, "-w", appath).Run(); err != nil {
			message = err.Error()
		}
		log.Warnf("%s: %s", tool, message)
	}
}

func exampleProjectPath(gopath string) (string, error) {
	var paths = strings.Split(gopath, ":")
	for _, p := range paths {
		p = path.Join(p, "src", examplePath)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("project template not found: %s", examplePath)
}

func newProject(log logger.Logger) cli.Command {
	action := func(ctx *cli.Context) error {
		var (
			ok         bool
			args       = ctx.Args()
			apprelpath string
			expath     string
		)

		if len(args) != 1 {
			log.Fatal("project path not set")
		}

		appath := args[0]
		appath, err := filepath.Abs(appath)
		if err != nil {
			log.Fatalf("project path is wrong: %v", err)
		}

		gopath := os.Getenv("GOPATH")

		if expath, err = exampleProjectPath(gopath); err != nil {
			log.Fatalf("project template not found: %v", err)
		}

		log.Info("Try to check project path")

		if ok, apprelpath = checkGopath(appath, gopath); !ok {
			log.Fatalf("project path must be in GOPATH(%s)", gopath)
		}

		if file, err := os.Stat(appath); err != nil {
			if err = os.Mkdir(appath, 0700); err != nil {
				log.Fatalf("Can't create project path: %v", err)
			}

			log.Info("Project path created")
		} else if !file.IsDir() {
			log.Fatal("Project path not folder")
		} else {
			log.Info("Project path found")
		}

		if ok, err := isEmptyDir(appath); err != nil {
			log.Fatalf("Can't check project path: %v", err)
		} else if !ok {
			log.Fatal("Project path must be empty")
		}

		if err := filepath.Walk(expath, func(p string, i os.FileInfo, err error) error {
			if p == expath {
				return nil
			}

			relpath := strings.Replace(p, expath, "", -1)

			pp := path.Join(appath, relpath)
			if i.IsDir() {
				if _, err := os.Stat(pp); err != nil {
					if err = os.Mkdir(pp, i.Mode()); err != nil {
						log.Fatalf("Can't create path: %v", err)
					}

					log.Infof("Create dir  %s", path.Join(apprelpath, relpath))
				}

				return nil
			}

			if path.Ext(p) != ".go" {
				if err := copyFileContents(p, pp); err != nil {
					log.Fatalf("Can't copy file: %v", err)
				}
			} else if err := copyGoFile(apprelpath, p, pp); err != nil {
				log.Fatalf("Can't copy file(%s): %v", i.Name(), err)
			}

			log.Infof("Create file %s", path.Join(apprelpath, relpath))

			return nil
		}); err != nil {
			log.Fatalf("Can't copy example project: %v", err)
		}

		projectFormatter(log, appath)

		log.Print("Project created")
		log.Printf(postInfo, appath)

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
