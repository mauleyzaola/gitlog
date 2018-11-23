package outputs

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/outputs/templates"
)

type HTMLOutput struct{}

func NewHTMLOutput() *HTMLOutput {
	return &HTMLOutput{}
}

func (t *HTMLOutput) DisplayCommits(data []byte) error {
	// create output directory
	dir, err := t.createDir()
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(dir, "charts.js"), []byte(templates.JS_COMMITS), 0666); err != nil {
		glog.Exit(err)
	}

	fileName := filepath.Join(dir, "index.html")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	raw := &struct {
		Raw string
	}{
		string(data),
	}

	if err = t.parseFile(templates.HTML_BASE, file, raw); err != nil {
		return err
	}

	return t.openUrl(fileName)
}

func (t *HTMLOutput) createDir() (string, error) {
	ts := time.Now().Format("20060102150403")
	dir := filepath.Join(os.TempDir(), ts)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func (t *HTMLOutput) parseFile(input string, file io.Writer, data interface{}) error {
	tpl := template.New("")
	tpl, err := tpl.Parse(input)
	if err != nil {
		return err
	}
	return tpl.Execute(file, data)
}

func (t *HTMLOutput) openUrl(uri string) error {
	var cmdName string
	switch runtime.GOOS {
	case "darwin":
		cmdName = "open"
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		return exec.Command(uri).Run()
	default:
		return fmt.Errorf("unsupported OS:%s", runtime.GOOS)

	}
	return exec.Command(cmdName, uri).Run()
}
