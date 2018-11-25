package outputs

import (
	"encoding/json"
	"fmt"
)

type JsonOutput struct{}

func NewJsonOutput() *JsonOutput {
	return &JsonOutput{}
}

func (t *JsonOutput) DisplayCommits(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
