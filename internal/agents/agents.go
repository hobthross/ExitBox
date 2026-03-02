package agents

import (
	"github.com/cloud-exit/exitbox/internal/agent"
	"github.com/cloud-exit/exitbox/internal/agents/claude"
	"github.com/cloud-exit/exitbox/internal/agents/codex"
	"github.com/cloud-exit/exitbox/internal/agents/opencode"
)

func Get(name string) agent.AgentEntity {
	return registry[name]
}

// registry holds all agent implementations.
var registry = map[string]agent.AgentEntity{}

// Register adds an agent to the registry.
func Register(a agent.AgentEntity) {
	registry[a.Name()] = a
}

func init() {
	Register(&claude.Claude{})
	Register(&codex.Codex{})
	Register(&opencode.OpenCode{})
}
