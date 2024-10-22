package schema

type Nex struct {
	Files []File `hcl:"file,block"`
}
