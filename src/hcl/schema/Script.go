package schema

type Script struct {
	Type string `hcl:"type,label"`

	ScriptSource string `hcl:"src"`
}

func (s Script) Handler() ScriptHandler {
	return ScriptHandler{
		Script: s,
	}
}

type ScriptHandler struct {
	Script Script
}

func (s ScriptHandler) Run() {
	switch s.Script.Type {
	case "shell":
		runScript(s.Script.ScriptSource)
	case "sh":
		runScript("/usr/bin/env", "sh", s.Script.ScriptSource)
	case "bash":
		runScript("/usr/bin/env", "bash", s.Script.ScriptSource)
	case "zsh":
		runScript("/usr/bin/env", "zsh", s.Script.ScriptSource)
	case "fish":
		runScript("/usr/bin/env", "fish", s.Script.ScriptSource)
	case "expect":
		runScript("/usr/bin/env", "expect", s.Script.ScriptSource)
	default:
	}
}
