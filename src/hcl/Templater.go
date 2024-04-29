package hcl

import (
	"fmt"
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

	// Output the result
	fmt.Println("Original:", statement)
	fmt.Println("Replaced:", replaced)
	return replaced
}
