package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/commit"
)

func newCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Manage GitLab repository commits",
	}
	cmd.AddCommand(newCommitCreateCmd())
	return cmd
}

func newCommitCreateCmd() *cobra.Command {
	var project string
	var specFile string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a commit from a JSON spec file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return commit.Create(cfg, p, specFile)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&specFile, "file", "", "JSON spec file path (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
