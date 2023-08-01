package resources

type Resources struct {
	env *Env
}

func (r *Resources) GetEnv() *Env {
	return r.env
}

func InitResources() *Resources {
	env := initEnv()

	return &Resources{env}
}
