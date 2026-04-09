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
	var projectID int
	var specFile string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a commit from a JSON spec file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commit.Create(mustConfig(cmd), projectID, specFile)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().StringVar(&specFile, "file", "", "JSON spec file path (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
