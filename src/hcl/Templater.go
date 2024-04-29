package hcl

import (
	"strings"
)

const (
	User = "{{USER}}"
)

func Scan(doc Provisioner, statement string) string {
	statement = user(doc, statement)
	return statement
}

func user(doc Provisioner, statement string) string {
	replaced := strings.Replace(statement, User, doc.PersonalUser, -1)
	return replaced
}
