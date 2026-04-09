package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/issuescheck"
)

func newIssuesCheckCmd() *cobra.Command {
	var projectID int
	var interactive, listOnly bool
	var closeIDs string
	cmd := &cobra.Command{
		Use:   "issues-check",
		Short: "Interactive issue review and bulk close tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issuescheck.Run(mustConfig(cmd), projectID, interactive, listOnly, closeIDs)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 1, "Project ID")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Interactive review mode")
	cmd.Flags().BoolVar(&listOnly, "list", false, "List issues only")
	cmd.Flags().StringVar(&closeIDs, "close", "", "Comma-separated issue IDs to close")
	return cmd
}
