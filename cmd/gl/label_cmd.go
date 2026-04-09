package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/label"
)

func newLabelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage GitLab labels",
	}
	cmd.AddCommand(newLabelListCmd(), newLabelCreateCmd())
	return cmd
}

func newLabelListCmd() *cobra.Command {
	var projectID, groupID, page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			return label.List(mustConfig(cmd), projectID, groupID, page, perPage)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID")
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}

func newLabelCreateCmd() *cobra.Command {
	var projectID, groupID int
	var name, color, desc string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a label",
		RunE: func(cmd *cobra.Command, args []string) error {
			return label.Create(mustConfig(cmd), projectID, groupID, name, color, desc)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID")
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID")
	cmd.Flags().StringVar(&name, "name", "", "Label name (required)")
	cmd.Flags().StringVar(&color, "color", "#428BCA", "Label color (hex)")
	cmd.Flags().StringVar(&desc, "description", "", "Label description")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
