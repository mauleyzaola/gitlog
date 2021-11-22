package outputs

import (
	"encoding/json"
	"fmt"
)

type JSONOutput struct{}

func NewJSONOutput() *JSONOutput {
	return &JSONOutput{}
}

func (t *JSONOutput) DisplayCommits(fg *FileGenerator, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
