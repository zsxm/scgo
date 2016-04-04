package chttp

type ControlConfig struct {
	Project string
	Module  string
	Title   string
	Comment string
}

func ControlConfigConfig() *ControlConfig {
	e := &ControlConfig{}
	return e
}
