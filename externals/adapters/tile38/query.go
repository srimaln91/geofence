package tile38

type Command struct {
	IsWrite    bool
	QueryParts []interface{}
}

func (command *Command) AddParam(param string) {
	command.QueryParts = append(command.QueryParts, param)
}
