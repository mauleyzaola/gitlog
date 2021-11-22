package git

import "strings"

// Author - The person which created the commit
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (t *Author) TrimEmailChars() string {
	res := strings.TrimPrefix(t.Email, "<")
	return strings.TrimSuffix(res, ">")
}
