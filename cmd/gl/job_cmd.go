package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/job"
)

func newJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "Manage and monitor GitLab CI/CD jobs",
	}
	cmd.AddCommand(
		newJobAnalyzeCmd(),
		newJobLogsCmd(),
		newJobRetryCmd(),
		newJobStatusCmd(),
		newJobCancelCmd(),
		newJobTriggerCmd(),
		newJobTraceCmd(),
		newJobDetailsCmd(),
	)
	return cmd
}

func newJobAnalyzeCmd() *cobra.Command {
	var jobID int
	var project, url string
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze a job (status, logs for failed jobs)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Analyze(cfg, jobID, p, url)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&url, "url", "", "GitLab job URL (alternative to --job)")
	return cmd
}

func newJobLogsCmd() *cobra.Command {
	var jobID, tail int
	var project string
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show job logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Logs(cfg, jobID, p, tail)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&tail, "tail", 100, "Number of log lines to show")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobRetryCmd() *cobra.Command {
	var jobID int
	var project string
	cmd := &cobra.Command{
		Use:   "retry",
		Short: "Retry a failed job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Retry(cfg, jobID, p)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobStatusCmd() *cobra.Command {
	var jobID int
	var project string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get job status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Status(cfg, jobID, p)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobCancelCmd() *cobra.Command {
	var jobID int
	var project string
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a running job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Cancel(cfg, jobID, p)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobTriggerCmd() *cobra.Command {
	var jobID int
	var project string
	cmd := &cobra.Command{
		Use:   "trigger",
		Short: "Trigger a manual job",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Trigger(cfg, jobID, p)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobTraceCmd() *cobra.Command {
	var jobID int
	var project, outputFile string
	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Get full job trace/log",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Trace(cfg, jobID, p, outputFile)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&outputFile, "output", "", "Save trace to file (default: stdout)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newJobDetailsCmd() *cobra.Command {
	var jobID int
	var project string
	cmd := &cobra.Command{
		Use:   "details",
		Short: "Show detailed job information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil { return err }
			return job.Details(cfg, jobID, p)
		},
	}
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}
