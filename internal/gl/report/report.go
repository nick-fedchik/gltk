package report

import (
	"fmt"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func newClient(cfg *config.Config) (*glclient.Client, error) {
	token := cfg.Token
	client, err := glclient.NewClient(token, glclient.WithBaseURL(cfg.GitLabURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return client, nil
}

func getProjectID(project string) interface{} {
	var id int
	if _, err := fmt.Sscanf(project, "%d", &id); err == nil {
		return id
	}
	return project
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func Pipelines(cfg *config.Config, project, branch, status string, limit int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListProjectPipelinesOptions{
		OrderBy: glclient.Ptr("updated_at"),
		Sort:    glclient.Ptr("desc"),
		ListOptions: glclient.ListOptions{
			PerPage: int64(limit),
		},
	}

	if branch != "" {
		opts.Ref = glclient.Ptr(branch)
	}

	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	if len(pipelines) == 0 {
		fmt.Println("No pipelines found")
		return nil
	}

	fmt.Printf("Pipelines Report (Last %d):\n\n", len(pipelines))
	for _, p := range pipelines {
		statusIcon := "✓"
		if p.Status == "failed" {
			statusIcon = "✗"
		} else if p.Status == "running" {
			statusIcon = "→"
		} else if p.Status == "canceled" {
			statusIcon = "⊘"
		}

		if status != "" && p.Status != status {
			continue
		}

		fmt.Printf("%s  #%d  %s  [%s]  %s\n",
			statusIcon,
			p.ID,
			p.Ref,
			p.Status,
			p.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
	return nil
}

func Jobs(cfg *config.Config, project string, pipelineID int64) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListJobsOptions{
		ListOptions: glclient.ListOptions{
			PerPage: 50,
		},
	}

	jobs, _, err := client.Jobs.ListPipelineJobs(projectID, pipelineID, opts)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No jobs found")
		return nil
	}

	fmt.Printf("Pipeline #%d Jobs (%d total):\n\n", pipelineID, len(jobs))
	for _, job := range jobs {
		statusIcon := "✓"
		if job.Status == "failed" {
			statusIcon = "✗"
		} else if job.Status == "running" {
			statusIcon = "→"
		} else if job.Status == "skipped" {
			statusIcon = "⊘"
		}

		duration := ""
		if job.Duration > 0 {
			duration = fmt.Sprintf(" (%.0f sec)", job.Duration)
		}

		fmt.Printf("%s  #%d  %-30s  [%s]%s\n",
			statusIcon,
			job.ID,
			truncate(job.Name, 30),
			job.Status,
			duration,
		)
	}
	return nil
}

func Summary(cfg *config.Config, project string, days int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	proj, _, err := client.Projects.GetProject(projectID, &glclient.GetProjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	opts := &glclient.ListProjectPipelinesOptions{
		OrderBy: glclient.Ptr("updated_at"),
		Sort:    glclient.Ptr("desc"),
		ListOptions: glclient.ListOptions{
			PerPage: 100,
		},
	}

	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	successCount := 0
	failedCount := 0
	runningCount := 0
	canceledCount := 0

	for _, p := range pipelines {
		switch p.Status {
		case "success":
			successCount++
		case "failed":
			failedCount++
		case "running":
			runningCount++
		case "canceled":
			canceledCount++
		}
	}

	fmt.Printf("Project Report: %s\n", proj.Name)
	fmt.Printf("Description: %s\n", proj.Description)
	fmt.Printf("\nStatistics (Last %d days):\n", days)
	fmt.Printf("  Total Pipelines: %d\n", len(pipelines))
	fmt.Printf("  Success: %d ✓\n", successCount)
	fmt.Printf("  Failed: %d ✗\n", failedCount)
	fmt.Printf("  Running: %d →\n", runningCount)
	fmt.Printf("  Canceled: %d ⊘\n", canceledCount)

	if len(pipelines) > 0 {
		successRate := (float64(successCount) / float64(len(pipelines))) * 100
		fmt.Printf("  Success Rate: %.1f%%\n", successRate)
	}

	fmt.Printf("\nProject Details:\n")
	fmt.Printf("  URL: %s\n", proj.WebURL)
	fmt.Printf("  Stars: %d ⭐\n", proj.StarCount)
	fmt.Printf("  Created: %s\n", proj.CreatedAt.Format("2006-01-02"))
	fmt.Printf("  Last Activity: %s\n", proj.LastActivityAt.Format("2006-01-02 15:04"))
	return nil
}
