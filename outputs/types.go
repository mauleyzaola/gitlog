package outputs

type Output interface {
	DisplayCommits(interface{}) error
}
