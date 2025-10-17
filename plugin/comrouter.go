package plugin

import "fmt"

// CommandHandler is a function that handles a specific command
type CommandHandler func(params map[string]any) Response

// CommandRouter helps route commands to handlers
type CommandRouter struct {
	handlers map[string]CommandHandler
}

// NewCommandRouter creates a new command router
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registers a handler for a command name
func (r *CommandRouter) Register(commandName string, handler CommandHandler) {
	r.handlers[commandName] = handler
}

// GetRegisteredCommands returns a list of all registered command names
func (r *CommandRouter) GetRegisteredCommands() []string {
	commands := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		commands = append(commands, name)
	}
	return commands
}

// Handle routes a command to the appropriate handler
func (r *CommandRouter) Handle(cmd Command) Response {
	handler, ok := r.handlers[cmd.Name]
	if !ok {
		return ErrorResponseMsg(fmt.Sprintf("unknown command: %s", cmd.Name))
	}
	// Params are already validated by the wrapper/datasource before reaching here
	return handler(cmd.Params)
}

// HandleJSON reads command from input, routes it, and writes response
func (r *CommandRouter) HandleJSON() int32 {
	cmd, err := ReadCommand()
	if err != nil {
		return WriteResponse(ErrorResponse(err))
	}
	return WriteResponse(r.Handle(cmd))
}
