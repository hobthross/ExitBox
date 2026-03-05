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
	"context"
	"os"

	"github.com/cloud-exit/exitbox/internal/agents"
	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/container"
	"github.com/cloud-exit/exitbox/internal/image"
	"github.com/cloud-exit/exitbox/internal/ui"
	"github.com/spf13/cobra"
)

var rebuildWorkspace string

var rebuildCmd = &cobra.Command{
	Use:   "rebuild <agent|all>",
	Short: "Force rebuild of agent image(s)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		rt := container.Detect()
		if rt == nil {
			ui.Error("No container runtime found. Install Podman or Docker.")
		}

		image.Version = Version
		image.AutoUpdate = true // rebuild always checks for latest

		var agentNames []string
		if name == "all" {
			cfg := config.LoadOrDefault()
			for _, a := range agents.All() {
				if cfg.IsAgentEnabled(a.Name()) {
					agentNames = append(agentNames, a.Name())
				}
			}
			if len(agentNames) == 0 {
				ui.Error("No agents are enabled. Run 'exitbox setup' first.")
			}
		} else {
			agt := agents.Get(name)
			if agt == nil {
				ui.Errorf("Unknown agent: %s", name)
			}
			agentNames = []string{agt.Name()}
		}

		projectDir, _ := os.Getwd()
		ctx := context.Background()

		for _, agentName := range agentNames {
			agt := agents.Get(agentName)
			if agt == nil {
				ui.Errorf("Unknown agent: %s", agentName)
				continue
			}

			ui.Infof("Rebuilding %s container image...", agt.DisplayName())
			if err := image.BuildCore(ctx, rt, agentName, true); err != nil {
				ui.Errorf("Failed to rebuild %s core image: %v", agt.DisplayName(), err)
			}
			if err := image.BuildProject(ctx, rt, agentName, projectDir, rebuildWorkspace, true); err != nil {
				ui.Errorf("Failed to rebuild %s project image: %v", agt.DisplayName(), err)
			}
			ui.Successf("%s image rebuilt successfully", agt.DisplayName())
		}
	},
}

func init() {
	rebuildCmd.Flags().StringVarP(&rebuildWorkspace, "workspace", "w", "", "Rebuild image for a specific workspace")
	rootCmd.AddCommand(rebuildCmd)
}
