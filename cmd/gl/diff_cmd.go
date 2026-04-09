package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/diff"
)

func newDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "View GitLab MR diffs",
	}
	cmd.AddCommand(newDiffSummaryCmd())
	return cmd
}

func newDiffSummaryCmd() *cobra.Command {
	var project string
	var mr int
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Show MR diff summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return diff.Summary(mustConfig(cmd), project, mr)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&mr, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}
