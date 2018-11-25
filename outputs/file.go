package outputs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/golang/glog"
)

type FileGenerator struct{}

func NewFileGenerator() (*FileGenerator, error) {
	return &FileGenerator{}, nil
}

func (t *FileGenerator) createDir() (string, error) {
	ts := time.Now().Format("20060102150403")
	dir := filepath.Join(os.TempDir(), ts)
	if err := os.MkdirAll(filepath.Join(dir, "lib"), 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func (t *FileGenerator) parseFile(input string, file io.Writer, data interface{}) error {
	tpl := template.New("")
	tpl, err := tpl.Parse(input)
	if err != nil {
		return err
	}
	return tpl.Execute(file, data)
}

func (t *FileGenerator) templateBox() *packr.Box {
	return packr.New("templates", "./templates")
}

func (t *FileGenerator) genCommitFiles(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// create output directory
	dir, err := t.createDir()
	if err != nil {
		return "", err
	}

	box := t.templateBox()
	commits, err := box.Find("commits.js")
	if err != nil {
		return "", err
	}

	if err = ioutil.WriteFile(filepath.Join(dir, "charts.js"), commits, 0666); err != nil {
		return "", err
	}

	// copy external libraries
	if err = box.Walk(func(name string, file packr.File) error {
		if strings.HasPrefix(name, "lib/") {
			return ioutil.WriteFile(filepath.Join(dir, name), []byte(file.String()), 0600)
		}
		return nil
	}); err != nil {
		return "", err
	}

	fileName := filepath.Join(dir, "index.html")
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer func() {
		if err = file.Close(); err != nil {
			glog.Error(err)
		}
	}()

	raw := &struct {
		Raw string
	}{
		string(data),
	}

	base, err := box.FindString("base.html")
	if err != nil {
		return "", err
	}

	if err = t.parseFile(base, file, raw); err != nil {
		return "", err
	}
	return dir, nil
}
