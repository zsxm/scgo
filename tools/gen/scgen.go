package gen

type GenEntity struct {
	GoPath, ProjectDir, GoSourceDir, ModuleName, FileDir string
}

const (
	GEN_ACTION  = "action"
	GEN_LOG     = "log"
	GEN_API     = "api"
	GEN_SERVICE = "service"
)
