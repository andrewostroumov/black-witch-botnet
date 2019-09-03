package relations

const (
	EventTypeHello = iota
	EventTypeRestart
)

const (
	ShellTypeExec = iota
	ShellTypeChangeDir
)

const (
	ErrorTimeout = iota
	ErrorCommand
	ErrorChangeDir
	ErrorUnknownRequest
	ErrorUnknownShellType
	ErrorUnknownEventType
)

type ShellCommand struct {
	Type uint8
	Data []byte
}

type EventMessage struct {
	Type uint8
	Data []byte
}

type ShellResult struct {
	Exit   int
	Stderr []byte
	Stdout []byte
}

type EventResult struct {
	Status bool
	Data   []byte
}

type ErrorResult struct {
	Code uint8
	Data []byte
}
