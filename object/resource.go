package object

type Resource struct {
	Name   string
	Handle any
}

func (r *Resource) Type() Type {
	return RESOURCE_OBJ
}

func (r *Resource) Inspect() string {
	return "<resource:" + r.Name + ">"
}
