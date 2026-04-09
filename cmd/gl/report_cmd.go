package main

import (
	"github.com/gltk/gltk/internal/gl/report"
	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate GitLab pipeline and job reports",
	}
	cmd.AddCommand(newReportPipelinesCmd(), newReportJobsCmd(), newReportSummaryCmd())
	return cmd
}

func newReportPipelinesCmd() *cobra.Command {
	var project, branch, status string
	var limit int
	cmd := &cobra.Command{
		Use:   "pipelines",
		Short: "Pipeline history report",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return report.Pipelines(cfg, p, branch, status, limit)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&branch, "branch", "", "Filter by branch")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status: success, failed, running, canceled")
	cmd.Flags().IntVar(&limit, "limit", 20, "Number of pipelines to show")
	return cmd
}

func newReportJobsCmd() *cobra.Command {
	var project string
	var pipelineID int64
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Job details report for a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return report.Jobs(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().Int64Var(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newReportSummaryCmd() *cobra.Command {
	var project string
	var days int
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Pipeline success/failure summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return report.Summary(cfg, p, days)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&days, "days", 7, "Number of days to include")
	return cmd
}
