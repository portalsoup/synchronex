package schema

import (
	"fmt"
	"log"
	"synchronex/common/execution"
)

type File struct {
	execution.Job
	Destination string `hcl:"type,label"`

	// If this file is to be copied, then it must have a source
	Source string `hcl:"src,optional"`

	User        string `hcl:"owner,optional"`
	Group       string `hcl:"group,optional"`
	Permissions string `hcl:"permissions,optional"`
}

func (f File) validation() (bool, error) {
	// Verify source path as file

	// verify containing folder of destination file can be read
	return false, nil
}

func (f File) execution() error {

	// Copy file from source to destination

	// chown to correct user/group

	//chmod to correct permissions like '0755'
	return nil
}

func (f File) ToString() string {
	_, err := f.validation()
	if err != nil {
		log.Fatalln(err)
	}
	msg := `
[%s]: %s
`
	return fmt.Sprintf(msg)

}
