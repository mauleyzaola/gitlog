package outputs

type Output interface {
	DisplayCommits(*FileGenerator, interface{}) error
}
