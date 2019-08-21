package relations

var Types = []string{"client", "shell"}
var Domains = []string{"exec", "cd"}

type Message struct {
	Type string
	Domain string
	Data string
}

type Result struct {
	Error string
	Exit uint8
	Data string
}
