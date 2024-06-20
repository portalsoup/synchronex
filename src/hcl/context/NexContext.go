package context

type NexContext struct {
	RequireRoot  bool   `hcl:"require_root,optional"`
	PersonalUser string `hcl:"user"`
}
