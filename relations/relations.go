package relations

var Types = []string{"client", "shell"}
var Domains = []string{"exec", "cd"}

type Command struct {
	Target string
	Scope  string
	Data   string
}

type Response struct {
	Error  *Error
	Result *Result
}

type Error struct {
	Code uint
	Data string
}

type Result struct {
	Exit   int
	Stderr []byte
	Stdout []byte
}
