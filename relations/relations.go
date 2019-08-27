package relations

var Types = []string{"payload", "shell"}
var Scopes = []string{"exec", "cd"}

const (
	TypeErrorResult = iota
	TypeSystemResult
	TypeShellResult
)

type Command struct {
	Target string
	Scope  string
	Data   string
}

type Response struct {
	Type uint8
	Data interface{}
}

type ErrorResult struct {
	Code uint
	Data string
}

type SystemResult struct {
	Status bool
}

type ShellResult struct {
	Exit   int
	Stderr []byte
	Stdout []byte
}
