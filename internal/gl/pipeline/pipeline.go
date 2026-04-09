package pipeline

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func resolveProjectID(client *glclient.Client, projectStr string) (int, error) {
	if id, err := strconv.Atoi(projectStr); err == nil {
		return id, nil
	}
	project, _, err := client.Projects.GetProject(projectStr, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve project %q: %w", projectStr, err)
	}
	return int(project.ID), nil
}

func List(cfg *config.Config, projectStr, status string, page int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	opts := &glclient.ListProjectPipelinesOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: 20,
		},
	}

	if status != "" {
		s := glclient.BuildStateValue(status)
		opts.Status = &s
	}

	pipelines, resp, err := client.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	for _, pipeline := range pipelines {
		createdAt := ""
		if pipeline.CreatedAt != nil {
			createdAt = pipeline.CreatedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("#%d: %s (%s) - %s\n",
			pipeline.ID,
			pipeline.Ref,
			pipeline.Status,
			createdAt)
	}

	if resp.NextPage > 0 {
		fmt.Printf("\nPage %d/%d | Next: %d\n", page, resp.TotalPages, resp.NextPage)
	}
	return nil
}

func Jobs(cfg *config.Config, projectStr string, pipelineID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	jobs, _, err := client.Jobs.ListPipelineJobs(projectID, int64(pipelineID), nil)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	for _, job := range jobs {
		duration := ""
		if job.Duration > 0 {
			duration = fmt.Sprintf(" (%.0f sec)", job.Duration)
		}
		fmt.Printf("  %d: %s - %s%s\n",
			job.ID,
			job.Name,
			job.Status,
			duration)
	}
	return nil
}

func Watch(cfg *config.Config, projectStr string, pipelineID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		p, _, err := client.Pipelines.GetPipeline(projectID, int64(pipelineID))
		if err != nil {
			return fmt.Errorf("failed to get pipeline: %w", err)
		}

		fmt.Print("\033[2J\033[H")
		fmt.Printf("Pipeline #%d Status: %s\n", p.ID, p.Status)
		fmt.Printf("Updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

		jobs, _, err := client.Jobs.ListPipelineJobs(projectID, int64(pipelineID), nil)
		if err == nil {
			for _, job := range jobs {
				symbol := "⏳"
				switch job.Status {
				case "success":
					symbol = "✅"
				case "failed":
					symbol = "❌"
				case "cancelled":
					symbol = "⏹️"
				case "skipped":
					symbol = "⏭️"
				}
				fmt.Printf("%s %s - %s\n", symbol, job.Name, job.Status)
			}
		}

		switch p.Status {
		case "success", "failed", "cancelled":
			fmt.Printf("\n✨ Pipeline finished with status: %s\n", p.Status)
			return nil
		}

		<-ticker.C
	}
}

func Trace(cfg *config.Config, projectStr string, jobID int, output string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	job, _, err := client.Jobs.GetJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Job #%d: %s | Stage: %s | Status: %s\n",
		job.ID, job.Name, job.Stage, job.Status)
	if job.Duration > 0 {
		fmt.Fprintf(os.Stderr, "Duration: %.0f seconds\n", job.Duration)
	}
	fmt.Fprintln(os.Stderr, "---")

	trace, _, err := client.Jobs.GetTraceFile(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("failed to get trace: %w", err)
	}

	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		if _, err := io.Copy(f, trace); err != nil {
			return fmt.Errorf("failed to write trace: %w", err)
		}
		fmt.Printf("✅ Trace saved to: %s\n", output)
	} else {
		if _, err := io.Copy(os.Stdout, trace); err != nil {
			return fmt.Errorf("failed to write trace: %w", err)
		}
	}
	return nil
}

func Cancel(cfg *config.Config, projectStr string, pipelineID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	p, _, err := client.Pipelines.CancelPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return fmt.Errorf("failed to cancel pipeline: %w", err)
	}

	fmt.Printf("⏹️  Pipeline #%d cancelled (status: %s)\n", p.ID, p.Status)
	return nil
}

func Create(cfg *config.Config, projectStr, ref, varsStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	opts := &glclient.CreatePipelineOptions{
		Ref: glclient.Ptr(ref),
	}

	p, _, err := client.Pipelines.CreatePipeline(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	fmt.Printf("✅ Pipeline #%d created successfully\n", p.ID)
	fmt.Printf("   Ref: %s\n", ref)
	fmt.Printf("   Status: %s\n", p.Status)
	fmt.Printf("   URL: %s/-/pipelines/%d\n", cfg.GitLabURL, p.ID)

	if varsStr != "" {
		fmt.Printf("   Variables: %s\n", varsStr)
	}
	return nil
}

func TriggerJob(cfg *config.Config, projectStr string, pipelineID int, jobName string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	jobOpts := &glclient.ListJobsOptions{
		ListOptions: glclient.ListOptions{
			PerPage: 100,
		},
	}

	jobs, _, err := client.Jobs.ListProjectJobs(projectID, jobOpts)
	if err != nil {
		return fmt.Errorf("failed to list pipeline jobs: %w", err)
	}

	var targetJob *glclient.Job
	for _, job := range jobs {
		if job.Name == jobName && job.Pipeline.ID == int64(pipelineID) {
			targetJob = job
			break
		}
	}

	if targetJob == nil {
		return fmt.Errorf("job '%s' not found in pipeline #%d", jobName, pipelineID)
	}

	triggeredJob, _, err := client.Jobs.PlayJob(projectID, targetJob.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to trigger job: %w", err)
	}

	fmt.Printf("▶️  Job '%s' (#%d) triggered successfully\n", jobName, triggeredJob.ID)
	fmt.Printf("   Status: %s\n", triggeredJob.Status)
	fmt.Printf("   Pipeline: #%d\n", pipelineID)
	fmt.Printf("   URL: %s/-/jobs/%d\n", cfg.GitLabURL, triggeredJob.ID)
	return nil
}

func TestReport(cfg *config.Config, projectStr string, pipelineID int, failedOnly bool) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	report, _, err := client.Pipelines.GetPipelineTestReport(projectID, int64(pipelineID))
	if err != nil {
		return fmt.Errorf("failed to get test report: %w", err)
	}

	if report.TotalCount == 0 {
		fmt.Println("No test results found for this pipeline.")
		fmt.Println("Hint: add artifacts:reports:junit to your .gitlab-ci.yml")
		return nil
	}

	// Summary
	fmt.Printf("Pipeline #%d — Test Report\n", pipelineID)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  Total: %d  ", report.TotalCount)
	fmt.Printf("✅ %d  ", report.SuccessCount)
	fmt.Printf("❌ %d  ", report.FailedCount)
	fmt.Printf("⏭ %d  ", report.SkippedCount)
	fmt.Printf("⚠ %d\n", report.ErrorCount)
	fmt.Printf("  Time:  %.2fs\n", report.TotalTime)
	fmt.Println()

	for _, suite := range report.TestSuites {
		if failedOnly && suite.FailedCount == 0 && suite.ErrorCount == 0 {
			continue
		}

		suiteStatus := "✅"
		if suite.FailedCount > 0 || suite.ErrorCount > 0 {
			suiteStatus = "❌"
		}
		fmt.Printf("%s %s (%d tests, %.2fs)\n", suiteStatus, suite.Name, suite.TotalCount, suite.TotalTime)

		for _, tc := range suite.TestCases {
			if failedOnly && tc.Status != "failed" && tc.Status != "error" {
				continue
			}

			symbol := "  ✅"
			switch tc.Status {
			case "failed":
				symbol = "  ❌"
			case "error":
				symbol = "  ⚠ "
			case "skipped":
				symbol = "  ⏭"
			}

			name := tc.Name
			if tc.Classname != "" && tc.Classname != tc.Name {
				name = tc.Classname + " > " + tc.Name
			}

			timeStr := ""
			if tc.ExecutionTime > 0 {
				timeStr = fmt.Sprintf(" (%.2fs)", tc.ExecutionTime)
			}

			fmt.Printf("%s %s%s\n", symbol, name, timeStr)

			if tc.StackTrace != "" {
				lines := splitLines(tc.StackTrace)
				limit := 10
				if len(lines) < limit {
					limit = len(lines)
				}
				for _, line := range lines[:limit] {
					fmt.Printf("       %s\n", line)
				}
				if len(lines) > 10 {
					fmt.Printf("       ... (%d more lines)\n", len(lines)-10)
				}
			}

			if tc.RecentFailures != nil && tc.RecentFailures.Count > 1 {
				fmt.Printf("       ⚡ Failed %d times on %s\n", tc.RecentFailures.Count, tc.RecentFailures.BaseBranch)
			}
		}
		fmt.Println()
	}

	return nil
}

func TestReportSummary(cfg *config.Config, projectStr string, pipelineID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	summary, _, err := client.Pipelines.GetPipelineTestReportSummary(projectID, int64(pipelineID))
	if err != nil {
		return fmt.Errorf("failed to get test report summary: %w", err)
	}

	if summary.Total.Count == 0 {
		fmt.Println("No test results found for this pipeline.")
		return nil
	}

	fmt.Printf("Pipeline #%d — Test Summary\n", pipelineID)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  Total: %d  ✅ %d  ❌ %d  ⏭ %d  ⚠ %d  (%.2fs)\n",
		summary.Total.Count, summary.Total.Success, summary.Total.Failed,
		summary.Total.Skipped, summary.Total.Error, summary.Total.Time)
	fmt.Println()

	for _, suite := range summary.TestSuites {
		status := "✅"
		if suite.FailedCount > 0 || suite.ErrorCount > 0 {
			status = "❌"
		}
		fmt.Printf("  %s %-40s  %3d total  %3d pass  %3d fail  (%.2fs)\n",
			status, suite.Name, suite.TotalCount, suite.SuccessCount, suite.FailedCount, suite.TotalTime)
		if suite.SuiteError != nil && *suite.SuiteError != "" {
			fmt.Printf("     Error: %s\n", *suite.SuiteError)
		}
	}

	return nil
}

func Coverage(cfg *config.Config, projectStr string, pipelineID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	p, _, err := client.Pipelines.GetPipeline(projectID, int64(pipelineID))
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %w", err)
	}

	fmt.Printf("Pipeline #%d — Coverage\n", pipelineID)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  Ref:    %s\n", p.Ref)
	fmt.Printf("  Status: %s\n", p.Status)

	if p.Coverage != "" {
		fmt.Printf("  Coverage: %s%%\n", p.Coverage)
	} else {
		fmt.Printf("  Coverage: not reported\n")
		fmt.Printf("  Hint: configure coverage regex in project CI/CD settings\n")
	}
	fmt.Println()

	// Per-job coverage breakdown
	jobs, _, err := client.Jobs.ListPipelineJobs(projectID, int64(pipelineID), nil)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	hasJobCoverage := false
	for _, job := range jobs {
		if job.Coverage > 0 {
			hasJobCoverage = true
			break
		}
	}

	if hasJobCoverage {
		fmt.Println("  Per-job coverage:")
		for _, job := range jobs {
			if job.Coverage > 0 {
				bar := coverageBar(job.Coverage)
				fmt.Printf("    %-30s %6.2f%%  %s\n", job.Name, job.Coverage, bar)
			}
		}
	} else if p.Coverage == "" {
		fmt.Println("  No job-level coverage data found.")
	}

	return nil
}

func coverageBar(pct float64) string {
	const width = 20
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	bar := make([]byte, width)
	for i := range bar {
		if i < filled {
			bar[i] = '#'
		} else {
			bar[i] = '-'
		}
	}
	return "[" + string(bar) + "]"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
