package artifact

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gltk/gltk/internal/config"
)

func List(cfg *config.Config, projectID string, jobID, page, perPage int) error {
	// If jobID specified, get artifact for that job
	if jobID != 0 {
		url := fmt.Sprintf("%s/api/v4/projects/%s/jobs/%d", cfg.GitLabURL, projectID, jobID)
		resp, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		resp.Header.Set("PRIVATE-TOKEN", cfg.Token)

		client := &http.Client{}
		res, err := client.Do(resp)
		if err != nil {
			return fmt.Errorf("failed to get job: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("failed to get job: HTTP %d", res.StatusCode)
		}

		var job struct {
			ID        int                      `json:"id"`
			Name      string                   `json:"name"`
			Artifacts []map[string]interface{} `json:"artifacts"`
		}
		if err := json.NewDecoder(res.Body).Decode(&job); err != nil {
			return fmt.Errorf("failed to decode job: %w", err)
		}

		if len(job.Artifacts) == 0 {
			fmt.Printf("Job #%d has no artifacts\n", jobID)
			return nil
		}

		fmt.Printf("Artifacts for Job #%d:\n", jobID)
		for _, artifact := range job.Artifacts {
			if fileType, ok := artifact["file_type"]; ok {
				fmt.Printf("  File Type: %v\n", fileType)
			}
			if fileFormat, ok := artifact["file_format"]; ok {
				fmt.Printf("    Format: %v\n", fileFormat)
			}
		}
		return nil
	}

	// Otherwise list jobs with artifacts
	url := fmt.Sprintf("%s/api/v4/projects/%s/jobs?page=%d&per_page=%d&status=success", cfg.GitLabURL, projectID, page, perPage)
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp.Header.Set("PRIVATE-TOKEN", cfg.Token)

	httpClient := &http.Client{}
	res, err := httpClient.Do(resp)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to list jobs: HTTP %d", res.StatusCode)
	}

	var jobs []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jobs); err != nil {
		return fmt.Errorf("failed to decode jobs: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No jobs found")
		return nil
	}

	fmt.Println("Jobs with artifacts:")
	count := 0
	for _, jobData := range jobs {
		if artifacts, ok := jobData["artifacts"].([]interface{}); ok && len(artifacts) > 0 {
			count++
			jobID := jobData["id"]
			jobName := jobData["name"]
			jobStatus := jobData["status"]
			fmt.Printf("  Job #%v: %v (%v)\n", jobID, jobName, jobStatus)
			for _, artifact := range artifacts {
				if artifactMap, ok := artifact.(map[string]interface{}); ok {
					fileType := artifactMap["file_type"]
					fileFormat := artifactMap["file_format"]
					fmt.Printf("    - %v (%v)\n", fileType, fileFormat)
				}
			}
		}
	}

	if count == 0 {
		fmt.Println("No artifacts found in successful jobs")
	}
	return nil
}

func Download(cfg *config.Config, projectID string, jobID int, outputPath string) error {
	// Construct download URL
	downloadURL := fmt.Sprintf("%s/api/v4/projects/%s/jobs/%d/artifacts", cfg.GitLabURL, projectID, jobID)

	// Determine output filename
	if outputPath == "" {
		outputPath = fmt.Sprintf("artifacts-job-%d.zip", jobID)
	}

	// Create request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)

	// Execute request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download artifact: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to download artifact: HTTP %d", res.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy artifact to file
	written, err := io.Copy(outFile, res.Body)
	if err != nil {
		return fmt.Errorf("failed to save artifact: %w", err)
	}

	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("✓ Artifact downloaded\n")
	fmt.Printf("  Job: #%d\n", jobID)
	fmt.Printf("  Size: %.2f MB\n", float64(written)/1024/1024)
	fmt.Printf("  Saved to: %s\n", absPath)
	return nil
}

func Delete(cfg *config.Config, projectID string, jobID int) error {
	// Construct delete URL
	deleteURL := fmt.Sprintf("%s/api/v4/projects/%s/jobs/%d/artifacts", cfg.GitLabURL, projectID, jobID)

	// Create DELETE request
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)

	// Execute request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 204 && res.StatusCode != 200 {
		return fmt.Errorf("failed to delete artifact: HTTP %d", res.StatusCode)
	}

	fmt.Printf("✓ Artifacts for job #%d deleted\n", jobID)
	return nil
}
