package main

import (
	"github.com/gltk/gltk/internal/gl/pipeline"
	"github.com/spf13/cobra"
)

func newPipelineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Monitor and manage GitLab CI/CD pipelines",
	}
	cmd.AddCommand(
		newPipelineListCmd(),
		newPipelineJobsCmd(),
		newPipelineWatchCmd(),
		newPipelineTraceCmd(),
		newPipelineCancelCmd(),
		newPipelineCreateCmd(),
		newPipelineTriggerJobCmd(),
		newPipelineTestReportCmd(),
		newPipelineTestSummaryCmd(),
		newPipelineCoverageCmd(),
	)
	return cmd
}

func newPipelineListCmd() *cobra.Command {
	var project, status string
	var page int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pipelines",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.List(cfg, p, status, page)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status: running, success, failed, cancelled")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	return cmd
}

func newPipelineJobsCmd() *cobra.Command {
	var project string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "List jobs in a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Jobs(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newPipelineWatchCmd() *cobra.Command {
	var project string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch pipeline until completion",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Watch(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newPipelineTraceCmd() *cobra.Command {
	var project, output string
	var jobID int
	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Get job trace/log",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Trace(cfg, p, jobID, output)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&output, "output", "", "Save trace to file (default: stdout)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newPipelineCancelCmd() *cobra.Command {
	var project string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a running pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Cancel(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newPipelineCreateCmd() *cobra.Command {
	var project, ref, vars string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Create(cfg, p, ref, vars)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Git reference: branch/tag/commit")
	cmd.Flags().StringVar(&vars, "vars", "", "Variables as KEY=VALUE,KEY2=VALUE2")
	return cmd
}

func newPipelineTriggerJobCmd() *cobra.Command {
	var project, jobName string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "trigger-job",
		Short: "Trigger a manual job in a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.TriggerJob(cfg, p, pipelineID, jobName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	cmd.Flags().StringVar(&jobName, "job", "", "Job name to trigger (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newPipelineTestReportCmd() *cobra.Command {
	var project string
	var pipelineID int
	var failedOnly bool
	cmd := &cobra.Command{
		Use:   "test-report",
		Short: "Show test results for a pipeline (JUnit)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.TestReport(cfg, p, pipelineID, failedOnly)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	cmd.Flags().BoolVar(&failedOnly, "failed", false, "Show only failed/errored tests")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newPipelineTestSummaryCmd() *cobra.Command {
	var project string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "test-summary",
		Short: "Show test suite summary for a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.TestReportSummary(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}

func newPipelineCoverageCmd() *cobra.Command {
	var project string
	var pipelineID int
	cmd := &cobra.Command{
		Use:   "coverage",
		Short: "Show code coverage for a pipeline and its jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return pipeline.Coverage(cfg, p, pipelineID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&pipelineID, "pipeline", 0, "Pipeline ID (required)")
	_ = cmd.MarkFlagRequired("pipeline")
	return cmd
}
