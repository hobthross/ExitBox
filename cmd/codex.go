// ExitBox - Multi-Agent Container Sandbox
// Copyright (C) 2026 Cloud Exit B.V.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"strings"

	"github.com/cloud-exit/exitbox/internal/agents/codex"
	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/profile"
	"github.com/cloud-exit/exitbox/internal/ui"
	"github.com/spf13/cobra"
)

func newCodexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codex",
		Short: "Manage Codex-specific features",
	}
	cmd.AddCommand(newCodexAccountsCmd())
	return cmd
}

func newCodexAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manage multiple Codex logins inside a workspace",
		Long: `Store, add, and switch named Codex accounts inside the active workspace.

The active Codex login remains in the workspace's normal .codex config path.
Saved accounts are stored in sidecar slots under the same workspace profile.`,
	}

	cmd.AddCommand(newCodexAccountsListCmd())
	cmd.AddCommand(newCodexAccountsCurrentCmd())
	cmd.AddCommand(newCodexAccountsSaveCmd())
	cmd.AddCommand(newCodexAccountsAddCmd())
	cmd.AddCommand(newCodexAccountsSwitchCmd())
	return cmd
}

func newCodexAccountsListCmd() *cobra.Command {
	var workspace string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved Codex accounts for a workspace",
		Run: func(cmd *cobra.Command, args []string) {
			ws, manager := resolveCodexAccountManager(workspace)
			items, err := manager.List()
			if err != nil {
				ui.Errorf("Failed to list Codex accounts: %v", err)
			}

			fmt.Printf("Workspace: %s\n", ws)
			if len(items) == 0 {
				fmt.Println("(no Codex accounts saved yet)")
				return
			}
			for _, item := range items {
				flags := make([]string, 0, 3)
				if item.Current {
					flags = append(flags, "current")
				}
				if item.Previous {
					flags = append(flags, "previous")
				}
				if !item.Ready {
					flags = append(flags, "pending-login")
				}
				if len(flags) == 0 {
					fmt.Printf("  %s\n", item.Name)
					continue
				}
				fmt.Printf("  %s (%s)\n", item.Name, strings.Join(flags, ", "))
			}
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace name (default: active workspace)")
	return cmd
}

func newCodexAccountsCurrentCmd() *cobra.Command {
	var workspace string
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current and previous Codex account names",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ws, manager := resolveCodexAccountManager(workspace)
			state, err := manager.LoadState()
			if err != nil {
				ui.Errorf("Failed to read Codex account state: %v", err)
			}
			fmt.Printf("Workspace: %s\n", ws)
			if state.Current != "" {
				fmt.Printf("Current: %s\n", state.Current)
			} else {
				fmt.Println("Current: (unknown)")
			}
			if state.Previous != "" {
				fmt.Printf("Previous: %s\n", state.Previous)
			}
		},
	}
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace name (default: active workspace)")
	return cmd
}

func newCodexAccountsSaveCmd() *cobra.Command {
	var workspace string
	cmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save the active Codex login under a name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ws, manager := resolveCodexAccountManager(workspace)
			if err := manager.Save(args[0]); err != nil {
				ui.Errorf("Failed to save Codex account in workspace '%s': %v", ws, err)
			}
			ui.Successf("Saved active Codex login as '%s' in workspace '%s'", args[0], ws)
		},
	}
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace name (default: active workspace)")
	return cmd
}

func newCodexAccountsAddCmd() *cobra.Command {
	var workspace string
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create an empty Codex account slot and make it current",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ws, manager := resolveCodexAccountManager(workspace)
			if err := manager.Add(args[0]); err != nil {
				ui.Errorf("Failed to add Codex account in workspace '%s': %v", ws, err)
			}
			ui.Successf("Added Codex account '%s' in workspace '%s'", args[0], ws)
			fmt.Printf("Next: exitbox run codex --workspace %s -- login\n", ws)
			fmt.Println("Run that from your host shell, not from inside the Codex prompt.")
			fmt.Printf("Optional: exitbox codex accounts save %s --workspace %s\n", args[0], ws)
		},
	}
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace name (default: active workspace)")
	return cmd
}

func newCodexAccountsSwitchCmd() *cobra.Command {
	var workspace string
	cmd := &cobra.Command{
		Use:   "switch <name>",
		Short: "Switch the workspace's active Codex login",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ws, manager := resolveCodexAccountManager(workspace)
			if err := manager.Switch(args[0]); err != nil {
				ui.Errorf("Failed to switch Codex account in workspace '%s': %v", ws, err)
			}
			ui.Successf("Switched Codex account to '%s' in workspace '%s'", args[0], ws)
		},
	}
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace name (default: active workspace)")
	return cmd
}

func resolveCodexAccountManager(workspace string) (string, *codex.AccountManager) {
	cfg := config.LoadOrDefault()
	ws := resolveConfigWorkspace(cfg, workspace)
	root := profile.WorkspaceAgentDir(ws, "codex")
	return ws, codex.NewAccountManager(root)
}

func init() {
	rootCmd.AddCommand(newCodexCmd())
}
