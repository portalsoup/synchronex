package schema

type Script struct {
	Type string `hcl:"type,label"`

	ScriptSource string `hcl:"src"`
}
