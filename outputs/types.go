package outputs

type Output interface {
	DisplayCommits(data []byte)
}
