package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/milestone"
)

func newMilestoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "Manage GitLab milestones",
	}
	cmd.AddCommand(
		newMilestoneListCmd(),
		newMilestoneCreateCmd(),
		newMilestoneUpdateCmd(),
		newMilestoneDeleteCmd(),
	)
	return cmd
}

func newMilestoneListCmd() *cobra.Command {
	var groupID, page, perPage int
	var state string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List milestones",
		RunE: func(cmd *cobra.Command, args []string) error {
			return milestone.List(mustConfig(cmd), groupID, state, page, perPage)
		},
	}
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID (required)")
	cmd.Flags().StringVar(&state, "state", "active", "Milestone state: active, closed")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	_ = cmd.MarkFlagRequired("group")
	return cmd
}

func newMilestoneCreateCmd() *cobra.Command {
	var groupID int
	var title, description, startDate, dueDate string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a milestone",
		RunE: func(cmd *cobra.Command, args []string) error {
			return milestone.Create(mustConfig(cmd), groupID, title, description, startDate, dueDate)
		},
	}
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID (required)")
	cmd.Flags().StringVar(&title, "title", "", "Milestone title (required)")
	cmd.Flags().StringVar(&description, "description", "", "Description")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	_ = cmd.MarkFlagRequired("group")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newMilestoneUpdateCmd() *cobra.Command {
	var groupID, milestoneID int
	var title, description, startDate, dueDate, state string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a milestone",
		RunE: func(cmd *cobra.Command, args []string) error {
			return milestone.Update(mustConfig(cmd), groupID, milestoneID, title, description, startDate, dueDate, state)
		},
	}
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID (required)")
	cmd.Flags().IntVar(&milestoneID, "id", 0, "Milestone ID (required)")
	cmd.Flags().StringVar(&title, "title", "", "Milestone title")
	cmd.Flags().StringVar(&description, "description", "", "Description")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&state, "state", "", "State event: activate, close")
	_ = cmd.MarkFlagRequired("group")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newMilestoneDeleteCmd() *cobra.Command {
	var groupID, milestoneID int
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a milestone",
		RunE: func(cmd *cobra.Command, args []string) error {
			return milestone.Delete(mustConfig(cmd), groupID, milestoneID)
		},
	}
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID (required)")
	cmd.Flags().IntVar(&milestoneID, "id", 0, "Milestone ID (required)")
	_ = cmd.MarkFlagRequired("group")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
