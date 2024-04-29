package schema

type Document struct {
	ProvisionerBlock Provisioner `hcl:"provisioner,block"`
}
