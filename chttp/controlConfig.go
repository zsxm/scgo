package chttp

type ControlConfig struct {
	project string
	module  string
	title   string
	comment string
}

type ControlConfigInterface interface {
	SetProject(project string)
	SetModule(module string)
	SetTitle(title string)
	SetComment(comment string)
	Project() string
	Module() string
	Title() string
	Comment() string
}

func ControlConfigConfig() ControlConfigInterface {
	e := &ControlConfig{}
	return e
}
func (this *ControlConfig) SetComment(comment string) {
	this.comment = comment
}
func (this *ControlConfig) Comment() string {
	return this.comment
}
func (this *ControlConfig) SetProject(project string) {
	this.project = project
}

func (this *ControlConfig) SetModule(module string) {
	this.module = module

}

func (this *ControlConfig) SetTitle(title string) {
	this.title = title

}

func (this *ControlConfig) Project() string {
	return this.project
}

func (this *ControlConfig) Module() string {
	return this.module
}

func (this *ControlConfig) Title() string {
	return this.title
}
