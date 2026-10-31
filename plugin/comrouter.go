package plugin

import (
	"fmt"

	m "github.com/plusev-terminal/go-plugin-common/meta"
)

// CommandHandler is a function that handles a specific command
type CommandHandler func(params map[string]any) Response

// CommandRouter helps route commands to handlers
type CommandRouter struct {
	handlers map[string]CommandHandler
	// priorities holds the plugin-declared rate-limit priority class per
	// command name (m.Priority* constants). Absent means "not declared";
	// declared values are NOT validated here — they travel the meta export
	// verbatim and the host validates them (falling back to its canonical
	// command-name table on unknown values).
	priorities map[string]string
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

// RegisterWithPriority registers a handler for a command name and declares
// its rate-limit priority class (m.PriorityTrading, m.PriorityInteractive,
// or m.PriorityBackground). An empty priority is equivalent to Register.
// Re-registering a command replaces both its handler and its declaration.
func (r *CommandRouter) RegisterWithPriority(commandName string, handler CommandHandler, priority string) {
	r.handlers[commandName] = handler
	if priority == "" {
		delete(r.priorities, commandName)
		return
	}
	if r.priorities == nil {
		r.priorities = make(map[string]string)
	}
	r.priorities[commandName] = priority
}

// CommandPriorities returns the declared command priorities sorted by
// command name for deterministic meta output; nil when none are declared.
func (r *CommandRouter) CommandPriorities() []m.CommandInfo {
	if len(r.priorities) == 0 {
		return nil
	}
	infos := make([]m.CommandInfo, 0, len(r.priorities))
	for name, priority := range r.priorities {
		infos = append(infos, m.CommandInfo{Name: name, Priority: priority})
	}
	return m.SortCommandInfosByName(infos)
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
