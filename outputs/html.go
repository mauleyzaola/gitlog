package outputs

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

type HTMLOutput struct{}

func NewHTMLOutput() *HTMLOutput {
	return &HTMLOutput{}
}

func (t *HTMLOutput) DisplayCommits(fg *FileGenerator, v interface{}) error {
	dirName, err := fg.genCommitFiles(v)
	if err != nil {
		return err
	}

	return t.openUrl(filepath.Join(dirName, "index.html"))
}

func (t *HTMLOutput) openUrl(uri string) error {
	var cmdName string
	switch runtime.GOOS {
	case "darwin":
		cmdName = "open"
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	default:
		return fmt.Errorf("unsupported OS:%s", runtime.GOOS)

	}
	return exec.Command(cmdName, uri).Run()
}
