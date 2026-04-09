package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/runner"
)

func newRunnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Manage GitLab CI/CD runners",
	}
	cmd.AddCommand(
		newRunnerListCmd(),
		newRunnerGetCmd(),
		newRunnerStatusCmd(),
		newRunnerDeleteCmd(),
		newRunnerPauseCmd(),
		newRunnerResumeCmd(),
		newRunnerJobsCmd(),
		newRunnerUpdateCmd(),
	)
	return cmd
}

func newRunnerListCmd() *cobra.Command {
	var runnerType, status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List runners",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.List(mustConfig(cmd), runnerType, status)
		},
	}
	cmd.Flags().StringVar(&runnerType, "type", "", "Runner type: instance_type, group_type, project_type")
	cmd.Flags().StringVar(&status, "status", "", "Runner status: online, offline, stale, never_contacted")
	return cmd
}

func newRunnerGetCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get runner details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Get(mustConfig(cmd), id)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Runner ID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newRunnerStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show runner status summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Status(mustConfig(cmd))
		},
	}
}

func newRunnerDeleteCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Delete(mustConfig(cmd), id)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Runner ID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newRunnerPauseCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause a runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Pause(mustConfig(cmd), id, true)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Runner ID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newRunnerResumeCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume a paused runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Resume(mustConfig(cmd), id)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Runner ID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newRunnerJobsCmd() *cobra.Command {
	var runnerID, limit int
	var status string
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "List jobs assigned to a runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Jobs(mustConfig(cmd), runnerID, limit, status)
		},
	}
	cmd.Flags().IntVar(&runnerID, "id", 0, "Runner ID (required)")
	cmd.Flags().IntVar(&limit, "limit", 20, "Number of jobs to show")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newRunnerUpdateCmd() *cobra.Command {
	var runnerID int
	var description, tags string
	var paused bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update runner settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Update(mustConfig(cmd), runnerID, description, tags, paused)
		},
	}
	cmd.Flags().IntVar(&runnerID, "id", 0, "Runner ID (required)")
	cmd.Flags().StringVar(&description, "description", "", "Runner description")
	cmd.Flags().StringVar(&tags, "tags", "", "Comma-separated tags")
	cmd.Flags().BoolVar(&paused, "paused", false, "Set paused state")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
