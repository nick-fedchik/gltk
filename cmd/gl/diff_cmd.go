package main

import (
	"github.com/gltk/gltk/internal/gl/diff"
	"github.com/spf13/cobra"
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
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return diff.Summary(cfg, p, mr)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&mr, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}
