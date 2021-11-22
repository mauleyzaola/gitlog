package outputs

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

type ZipOutput struct {
	file string
}

func NewZipOutput(file string) (*ZipOutput, error) {
	if file == "" {
		return nil, errors.New("missing filename. cannot generate zip files from empty target")
	}
	return &ZipOutput{
		file: file,
	}, nil
}

func (t *ZipOutput) DisplayCommits(fg *FileGenerator, v interface{}) error {
	dir, err := fg.genCommitFiles(v)
	if err != nil {
		return err
	}
	zipFile, err := os.Create(t.file)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if err = zipWriter.Close(); err != nil {
			glog.Error(err)
		}
		if err = zipFile.Close(); err != nil {
			glog.Error(err)
		}
	}()

	if errFile := t.addFiles(zipWriter, dir, ""); errFile != nil {
		return errFile
	}
	return nil
}

func (t *ZipOutput) addFiles(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := os.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			return t.addFiles(w, filepath.Join(basePath, file.Name()), file.Name()+"/")
		}
	}
	return nil
}
