package cmd

import (
	"fmt"
	"os"

	"github.com/cloud-exit/exitbox/internal/agents"
	"github.com/cloud-exit/exitbox/internal/plugins"
	"github.com/cloud-exit/exitbox/internal/ui"
	"github.com/spf13/cobra"
)

func newPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage agent plugins",
		Long:  "Clone, list, and remove workspace-scoped plugins for supported agents.",
	}

	cmd.AddCommand(newPluginInstallCmd())
	cmd.AddCommand(newPluginListCmd())
	cmd.AddCommand(newPluginRemoveCmd())
	return cmd
}

func newPluginInstallCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "install <agent> <git-url-or-path>",
		Short: "Install a plugin by cloning its repository",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			agentName := args[0]
			source := args[1]

			if agents.Get(agentName) == nil {
				ui.Errorf("Unknown agent: %s", agentName)
			}

			projectDir, err := os.Getwd()
			if err != nil {
				ui.Errorf("Failed to determine current directory: %v", err)
			}
			installed, err := plugins.Install(projectDir, agentName, source, name)
			if err != nil {
				ui.Errorf("Failed to install plugin: %v", err)
			}
			if len(installed) == 1 {
				ui.Successf("Installed plugin '%s' for %s in project '%s'", installed[0].Name, agentName, projectDir)
				fmt.Printf("Host path: %s\n", installed[0].Dir)
				return
			}
			ui.Successf("Installed %d plugins for %s in project '%s'", len(installed), agentName, projectDir)
			for _, plugin := range installed {
				fmt.Printf("%s -> %s\n", plugin.Name, plugin.Dir)
			}
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Override plugin directory name")
	return cmd
}

func newPluginListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <agent>",
		Short: "List installed plugins for an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			agentName := args[0]
			if agents.Get(agentName) == nil {
				ui.Errorf("Unknown agent: %s", agentName)
			}

			projectDir, err := os.Getwd()
			if err != nil {
				ui.Errorf("Failed to determine current directory: %v", err)
			}
			list, err := plugins.List(projectDir, agentName)
			if err != nil {
				ui.Errorf("Failed to list plugins: %v", err)
			}
			if len(list) == 0 {
				fmt.Println("(no plugins installed)")
				return
			}
			for _, plugin := range list {
				fmt.Println(plugin.Name)
			}
		},
	}
	return cmd
}

func newPluginRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <agent> <name>",
		Short: "Remove an installed plugin",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			agentName := args[0]
			name := args[1]

			if agents.Get(agentName) == nil {
				ui.Errorf("Unknown agent: %s", agentName)
			}

			projectDir, err := os.Getwd()
			if err != nil {
				ui.Errorf("Failed to determine current directory: %v", err)
			}
			if err := plugins.Remove(projectDir, agentName, name); err != nil {
				ui.Errorf("Failed to remove plugin: %v", err)
			}
			ui.Successf("Removed plugin '%s' for %s from project '%s'", name, agentName, projectDir)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newPluginCmd())
}
