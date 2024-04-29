package hcl

import (
	"strings"
	"synchronex/src/hcl/schema"
)

const (
	User = "{{USER}}"
)

func Scan(doc schema.Provisioner, statement string) string {
	statement = user(doc, statement)
	return statement
}

func user(doc schema.Provisioner, statement string) string {
	replaced := strings.Replace(statement, User, doc.PersonalUser, -1)
	return replaced
}
