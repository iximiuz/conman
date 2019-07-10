package oci

type runcRuntime struct {
	exePath  string
	rootPath string
}

func NewRuntime(exePath, rootPath string) Runtime {
	return &runcRuntime{
		exePath:  exePath,
		rootPath: rootPath,
	}
}

func (r *runcRuntime) CreateContainer() {
	panic("not implemented")
}

func (r *runcRuntime) StartContainer() {
	panic("not implemented")
}

func (r *runcRuntime) KillContainer() {
	panic("not implemented")
}

func (r *runcRuntime) DeleteContainer() {
	panic("not implemented")
}

func (r *runcRuntime) ContainerState() {
	panic("not implemented")
}
