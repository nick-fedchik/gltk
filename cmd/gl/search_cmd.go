package main

import (
	"github.com/gltk/gltk/internal/gl/search"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search GitLab issues and merge requests",
	}
	cmd.AddCommand(newSearchIssuesCmd(), newSearchMRsCmd())
	return cmd
}

func newSearchIssuesCmd() *cobra.Command {
	var project, text, state string
	var page, perPage int
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "Search issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return search.Issues(cfg, p, text, state, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&text, "search", "", "Search text")
	cmd.Flags().StringVar(&state, "state", "opened", "State: opened, closed, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}

func newSearchMRsCmd() *cobra.Command {
	var project, text, state string
	var page, perPage int
	cmd := &cobra.Command{
		Use:   "mrs",
		Short: "Search merge requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return search.MRs(cfg, p, text, state, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&text, "search", "", "Search text")
	cmd.Flags().StringVar(&state, "state", "opened", "State: opened, closed, merged, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}
