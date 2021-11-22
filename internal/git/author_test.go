package git

import "testing"

func TestAuthor_TrimEmailChars(t *testing.T) {
	author := &Author{
		Name:  "mauleyzaola",
		Email: "<mauricio.leyzaola@gmail.com>",
	}
	author.Email = author.TrimEmailChars()
	if expected, actual := "mauricio.leyzaola@gmail.com", author.Email; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}
