package agents

import (
	"github.com/cloud-exit/exitbox/internal/agent"
	"github.com/cloud-exit/exitbox/internal/agents/claude"
	"github.com/cloud-exit/exitbox/internal/agents/codex"
	"github.com/cloud-exit/exitbox/internal/agents/opencode"
)

func Get(name string) agent.Agent {
	return registry[name]
}

func All() []agent.Agent {
	result := make([]agent.Agent, 0, len(registry))
	for _, a := range registry {
		result = append(result, a)
	}
	return result
}

// Names returns agent names from the registry. If prefix is passed, those strings
// are prepended to the result.
func Names(prefix ...string) []string {
	names := make([]string, 0, len(prefix)+len(registry))
	names = append(names, prefix...)
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// registry holds all agent implementations.
var registry = map[string]agent.Agent{}

// Register adds an agent to the registry.
func Register(a agent.Agent) {
	registry[a.Name()] = a
}

func init() {
	Register(&claude.Claude{})
	Register(&codex.Codex{})
	Register(&opencode.OpenCode{})
}
