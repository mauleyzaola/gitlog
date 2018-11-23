package outputs

import "fmt"

type JsonOutput struct {}


func NewJsonOutput()*JsonOutput{
	return &JsonOutput{}
}


func (t *JsonOutput) DisplayCommits(data []byte) {
	fmt.Print(string(data))
}
