package template

import (
	"strings"
)

const (
	User = "{{USER}}"
)

// ReplaceUser for mustache template variables
func ReplaceUser(user, statement string) string {
	replaced := strings.Replace(statement, User, user, -1)
	return replaced
}
