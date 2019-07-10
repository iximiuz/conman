package rollback

type Rollback struct {
	funcs []func()
}

func New() *Rollback {
	return &Rollback{}
}

func (r *Rollback) Add(fn func()) {
	r.funcs = append(r.funcs, fn)
}

func (r *Rollback) Execute() bool {
	for _, fn := range r.funcs {
		fn()
	}
	return true
}
