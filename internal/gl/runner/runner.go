package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gltk/gltk/internal/config"
)

func doRequest(method, url, token string, payload interface{}) ([]byte, error) {
	var reqBody io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GitLab API error %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

func printJSON(data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(string(data))
		return
	}
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

func strOrDash(v interface{}) string {
	if v == nil {
		return "-"
	}
	s, _ := v.(string)
	if s == "" {
		return "-"
	}
	return s
}

func intOrDash(v interface{}) string {
	if v == nil {
		return "-"
	}
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%d", int(n))
	case int:
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%v", v)
}

func boolStr(v interface{}) string {
	if b, ok := v.(bool); ok && b {
		return "yes"
	}
	return "no"
}

func List(cfg *config.Config, runnerType, status string) error {
	token := cfg.Token
	url := cfg.GitLabURL + "/api/v4/runners/all?per_page=100"
	if runnerType != "" {
		url += "&type=" + runnerType
	}
	if status != "" {
		url += "&status=" + status
	}

	body, err := doRequest("GET", url, token, nil)
	if err != nil {
		return err
	}

	var runners []map[string]interface{}
	if err := json.Unmarshal(body, &runners); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tONLINE\tPAUSED\tTAGS")
	fmt.Fprintln(w, "--\t----\t------\t------\t------\t----")
	for _, r := range runners {
		tags := ""
		if t, ok := r["tag_list"].([]interface{}); ok {
			for i, tag := range t {
				if i > 0 {
					tags += ","
				}
				tags += fmt.Sprintf("%v", tag)
			}
		}
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
			intOrDash(r["id"]),
			strOrDash(r["description"]),
			strOrDash(r["status"]),
			boolStr(r["online"]),
			boolStr(r["paused"]),
			tags,
		)
	}
	w.Flush()
	fmt.Printf("\nTotal: %d runners\n", len(runners))
	return nil
}

func Get(cfg *config.Config, id int) error {
	token := cfg.Token
	url := fmt.Sprintf("%s/api/v4/runners/%d", cfg.GitLabURL, id)
	body, err := doRequest("GET", url, token, nil)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func Status(cfg *config.Config) error {
	token := cfg.Token
	body, err := doRequest("GET", cfg.GitLabURL+"/api/v4/runners/all?per_page=100", token, nil)
	if err != nil {
		return err
	}

	var runners []map[string]interface{}
	if err := json.Unmarshal(body, &runners); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("GitLab Runners Status\n")
	fmt.Printf("══════════════════════════════════════════\n")

	online, offline, paused := 0, 0, 0
	for _, r := range runners {
		id := intOrDash(r["id"])
		name := strOrDash(r["description"])
		status := strOrDash(r["status"])
		isOnline := r["online"] == true
		isPaused := r["paused"] == true

		icon := "🔴"
		if isPaused {
			icon = "⏸️ "
			paused++
		} else if isOnline {
			icon = "🟢"
			online++
		} else {
			offline++
		}

		tags := ""
		if t, ok := r["tag_list"].([]interface{}); ok && len(t) > 0 {
			tags = " ["
			for i, tag := range t {
				if i > 0 {
					tags += ", "
				}
				tags += fmt.Sprintf("%v", tag)
			}
			tags += "]"
		}

		fmt.Printf("%s  ID=%-4s  %-30s  status=%-12s%s\n", icon, id, name, status, tags)
	}

	fmt.Printf("──────────────────────────────────────────\n")
	fmt.Printf("Total: %d  🟢 Online: %d  🔴 Offline: %d  ⏸️  Paused: %d\n",
		len(runners), online, offline, paused)

	if online == 0 && len(runners) > 0 {
		fmt.Println("\n⚠️  WARNING: No runners are online! CI/CD jobs will be stuck.")
	}
	return nil
}

func Delete(cfg *config.Config, id int) error {
	token := cfg.Token
	url := fmt.Sprintf("%s/api/v4/runners/%d", cfg.GitLabURL, id)
	_, err := doRequest("DELETE", url, token, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Runner id=%d deleted.\n", id)
	return nil
}

func Pause(cfg *config.Config, id int, pause bool) error {
	token := cfg.Token
	payload := map[string]interface{}{"paused": pause}
	url := fmt.Sprintf("%s/api/v4/runners/%d", cfg.GitLabURL, id)
	_, err := doRequest("PUT", url, token, payload)
	if err != nil {
		return err
	}

	action := "resumed"
	if pause {
		action = "paused"
	}
	fmt.Printf("Runner id=%d %s.\n", id, action)
	return nil
}

func Resume(cfg *config.Config, id int) error {
	return Pause(cfg, id, false)
}

func Jobs(cfg *config.Config, runnerID, limit int, status string) error {
	token := cfg.Token
	url := fmt.Sprintf("%s/api/v4/runners/%d/jobs?per_page=%d", cfg.GitLabURL, runnerID, limit)
	if status != "" {
		url += "&status=" + status
	}

	body, err := doRequest("GET", url, token, nil)
	if err != nil {
		return err
	}

	var jobs []map[string]interface{}
	if err := json.Unmarshal(body, &jobs); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No jobs found for this runner")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tPROJECT\tSTAGE\tDURATION")
	fmt.Fprintln(w, "--\t----\t------\t-------\t-----\t--------")

	for _, job := range jobs {
		id := intOrDash(job["id"])
		name := strOrDash(job["name"])
		jobStatus := strOrDash(job["status"])
		project := ""
		if p, ok := job["project"].(map[string]interface{}); ok {
			project = strOrDash(p["name"])
		}
		stage := strOrDash(job["stage"])
		duration := intOrDash(job["duration"])

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", id, name, jobStatus, project, stage, duration)
	}
	w.Flush()
	return nil
}

func Update(cfg *config.Config, runnerID int, description, tagsStr string, paused bool) error {
	token := cfg.Token

	payload := make(map[string]interface{})

	if description != "" {
		payload["description"] = description
	}

	if tagsStr != "" {
		var tags []string
		for _, tag := range strings.Split(tagsStr, ",") {
			t := strings.TrimSpace(tag)
			if t != "" {
				tags = append(tags, t)
			}
		}
		payload["tag_list"] = tags
	}

	if paused {
		payload["paused"] = true
	}

	if len(payload) == 0 {
		return fmt.Errorf("at least one update option required (--description, --tags, --paused)")
	}

	url := fmt.Sprintf("%s/api/v4/runners/%d", cfg.GitLabURL, runnerID)
	_, err := doRequest("PUT", url, token, payload)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Runner id=%d updated successfully\n", runnerID)
	if description != "" {
		fmt.Printf("   Description: %s\n", description)
	}
	if tagsStr != "" {
		fmt.Printf("   Tags: %s\n", tagsStr)
	}
	if paused {
		fmt.Printf("   Status: paused\n")
	}
	return nil
}
