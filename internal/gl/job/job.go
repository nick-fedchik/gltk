package job

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func resolveProjectID(client *glclient.Client, projectStr string) (int64, error) {
	if id, err := strconv.ParseInt(projectStr, 10, 64); err == nil {
		return id, nil
	}
	project, _, err := client.Projects.GetProject(projectStr, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve project %q: %w", projectStr, err)
	}
	return project.ID, nil
}

func readTrace(r io.Reader) string {
	if r == nil {
		return ""
	}
	data, _ := io.ReadAll(r)
	return string(data)
}

func Analyze(cfg *config.Config, jobID int, projectStr, urlFlag string) error {
	if urlFlag != "" {
		urlProject, urlJobID := extractFromURL(urlFlag)
		if urlJobID > 0 {
			jobID = urlJobID
		}
		if urlProject != "" && projectStr == "1" {
			projectStr = urlProject
		}
	}

	if jobID == 0 {
		return fmt.Errorf("--job or --url required")
	}

	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.GetJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error fetching job: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	fmt.Printf("\n📋 JOB ANALYSIS #%d\n", jobID)
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Name:       %s\n", j.Name)
	fmt.Printf("Status:     %s\n", colorStatus(j.Status))
	fmt.Printf("Stage:      %s\n", j.Stage)
	fmt.Printf("Duration:   %.0f seconds\n", j.Duration)
	fmt.Printf("Created:    %v\n", j.CreatedAt)
	fmt.Printf("Started:    %v\n", j.StartedAt)
	fmt.Printf("Finished:   %v\n", j.FinishedAt)

	if j.Status == "failed" {
		fmt.Println("\n🔴 JOB FAILED - Fetching error details...")
		fmt.Println("─────────────────────────────────────")
		trace, _, err := client.Jobs.GetTraceFile(projectID, int64(jobID))
		if err == nil && trace != nil {
			traceStr := readTrace(trace)
			lines := strings.Split(traceStr, "\n")
			start := len(lines) - 30
			if start < 0 {
				start = 0
			}
			for _, line := range lines[start:] {
				if line != "" {
					fmt.Println(line)
				}
			}
		}
		fmt.Println("\n💡 RECOMMENDATIONS:")
		fmt.Println("  • Run: ./bin/gl-job logs --job=" + fmt.Sprintf("%d", jobID) + " --tail=100")
		fmt.Println("  • Fix the issue and retry: ./bin/gl-job retry --job=" + fmt.Sprintf("%d", jobID))
	} else {
		fmt.Printf("\n✅ Job Status: %s\n", j.Status)
	}
	return nil
}

func Logs(cfg *config.Config, jobID int, projectStr string, tail int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	trace, _, err := client.Jobs.GetTraceFile(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error fetching logs: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	traceStr := readTrace(trace)
	lines := strings.Split(traceStr, "\n")
	start := len(lines) - tail
	if start < 0 {
		start = 0
	}

	fmt.Printf("📄 Last %d lines of job #%d:\n", tail, jobID)
	fmt.Println("═════════════════════════════════════════")
	for _, line := range lines[start:] {
		fmt.Println(line)
	}
	return nil
}

func Retry(cfg *config.Config, jobID int, projectStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.RetryJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error retrying job: %w", err)
	}

	fmt.Printf("✅ Job #%d restarted\n", j.ID)
	fmt.Printf("   Status: %s\n", j.Status)
	fmt.Printf("   Pipeline: %s/-/pipelines/%d\n", cfg.GitLabURL, j.Pipeline.ID)
	return nil
}

func Status(cfg *config.Config, jobID int, projectStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.GetJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error fetching job: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	fmt.Printf("Job #%d: %s (%s)\n", j.ID, j.Name, colorStatus(j.Status))
	fmt.Printf("  Stage:    %s\n", j.Stage)
	fmt.Printf("  Duration: %.0f seconds\n", j.Duration)
	return nil
}

func Cancel(cfg *config.Config, jobID int, projectStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.CancelJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error cancelling job: %w", err)
	}

	fmt.Printf("⏹️  Job #%d cancelled (status: %s)\n", j.ID, j.Status)
	return nil
}

func Trigger(cfg *config.Config, jobID int, projectStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.PlayJob(projectID, int64(jobID), nil)
	if err != nil {
		return fmt.Errorf("error triggering job: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	fmt.Printf("▶️  Job #%d triggered successfully\n", j.ID)
	fmt.Printf("   Name: %s\n", j.Name)
	fmt.Printf("   Status: %s\n", colorStatus(j.Status))
	fmt.Printf("   Pipeline: %s/-/pipelines/%d\n", cfg.GitLabURL, j.Pipeline.ID)
	return nil
}

func Trace(cfg *config.Config, jobID int, projectStr, outputFile string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	trace, _, err := client.Jobs.GetTraceFile(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error fetching trace: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	traceStr := readTrace(trace)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(traceStr), 0644); err != nil {
			return fmt.Errorf("error saving to file: %w", err)
		}
		fmt.Printf("✅ Full trace saved to %s\n", outputFile)
		fmt.Printf("   Lines: %d\n", len(strings.Split(traceStr, "\n")))
		return nil
	}

	lines := strings.Split(traceStr, "\n")
	fmt.Printf("📋 Full trace for job #%d (%d lines):\n", jobID, len(lines))
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(traceStr)
	return nil
}

func Details(cfg *config.Config, jobID int, projectStr string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projectID, err := resolveProjectID(client, projectStr)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}

	j, _, err := client.Jobs.GetJob(projectID, int64(jobID))
	if err != nil {
		return fmt.Errorf("error fetching job: %w (project=%s, job=%d)", err, projectStr, jobID)
	}

	type JobDetails struct {
		ID           int64       `json:"id"`
		Name         string      `json:"name"`
		Status       string      `json:"status"`
		Stage        string      `json:"stage"`
		CreatedAt    interface{} `json:"created_at"`
		StartedAt    interface{} `json:"started_at"`
		FinishedAt   interface{} `json:"finished_at"`
		Duration     float64     `json:"duration"`
		WebURL       string      `json:"web_url"`
		PipelineID   int64       `json:"pipeline_id"`
		Runner       interface{} `json:"runner"`
		Tags         []string    `json:"tags"`
		AllowFailure bool        `json:"allow_failure"`
	}

	details := JobDetails{
		ID:           j.ID,
		Name:         j.Name,
		Status:       j.Status,
		Stage:        j.Stage,
		CreatedAt:    j.CreatedAt,
		StartedAt:    j.StartedAt,
		FinishedAt:   j.FinishedAt,
		Duration:     j.Duration,
		WebURL:       j.WebURL,
		PipelineID:   j.Pipeline.ID,
		Runner:       j.Runner,
		Tags:         j.TagList,
		AllowFailure: j.AllowFailure,
	}

	jsonData, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

func extractFromURL(urlStr string) (projectPath string, jobID int) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", 0
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	if idx := strings.Index(path, "/-/"); idx >= 0 {
		projectPath = path[:idx]
		rest := path[idx+3:]
		parts := strings.Split(rest, "/")
		if len(parts) >= 2 && parts[0] == "jobs" {
			jobID, _ = strconv.Atoi(parts[1])
		}
	}

	if jobID == 0 {
		re := regexp.MustCompile(`jobs/(\d+)`)
		matches := re.FindStringSubmatch(urlStr)
		if len(matches) > 1 {
			jobID, _ = strconv.Atoi(matches[1])
		}
	}

	return projectPath, jobID
}

func colorStatus(status string) string {
	switch status {
	case "success":
		return "✅ success"
	case "failed":
		return "❌ failed"
	case "running":
		return "🔵 running"
	case "pending":
		return "⏳ pending"
	case "canceled":
		return "⚫ canceled"
	default:
		return status
	}
}
