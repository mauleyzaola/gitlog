package outputs

import (
	"embed"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/golang/glog"
)

//go:embed templates
var templateFiles embed.FS

type FileGenerator struct{}

func NewFileGenerator() (*FileGenerator, error) {
	return &FileGenerator{}, nil
}

func (t *FileGenerator) createDir() (string, error) {
	ts := time.Now().Format("20060102150403")
	dir := filepath.Join(os.TempDir(), ts)
	if err := os.MkdirAll(filepath.Join(dir, "lib"), os.ModePerm); err != nil {
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

func (t *FileGenerator) readFile(name string) ([]byte, error) {
	file, err := templateFiles.Open(name)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return io.ReadAll(file)
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

	commits, err := t.readFile(filepath.Join("templates", "commits.js"))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	if err = ioutil.WriteFile(filepath.Join(dir, "charts.js"), commits, 0666); err != nil {
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

	base, err := t.readFile(filepath.Join("templates", "base.html"))
	if err != nil {
		glog.Error(err)
		return "", err
	}

	if err = t.parseFile(string(base), file, raw); err != nil {
		return "", err
	}
	return dir, nil
}
