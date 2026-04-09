package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func resolveUsername(baseURL, token, username string) (int, error) {
	url := fmt.Sprintf("%s/api/v4/users?username=%s", baseURL, username)
	body, err := doRequest("GET", url, token, nil)
	if err != nil {
		return 0, err
	}
	var users []map[string]interface{}
	if err := json.Unmarshal(body, &users); err != nil || len(users) == 0 {
		return 0, fmt.Errorf("user '%s' not found", username)
	}
	id, ok := users[0]["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected id type for user '%s'", username)
	}
	return int(id), nil
}

func List(cfg *config.Config, search string, active bool, page, perPage int) error {
	token := cfg.Token
	url := fmt.Sprintf("%s/api/v4/users?page=%d&per_page=%d", cfg.GitLabURL, page, perPage)
	if search != "" {
		url += "&search=" + search
	}
	if active {
		url += "&active=true"
	}

	body, err := doRequest("GET", url, token, nil)
	if err != nil {
		return err
	}

	var users []map[string]interface{}
	if err := json.Unmarshal(body, &users); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUSERNAME\tNAME\tSTATE\tADMIN\tEMAIL")
	fmt.Fprintln(w, "--\t--------\t----\t-----\t-----\t-----")
	for _, u := range users {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
			intOrDash(u["id"]),
			strOrDash(u["username"]),
			strOrDash(u["name"]),
			strOrDash(u["state"]),
			boolStr(u["is_admin"]),
			strOrDash(u["email"]),
		)
	}
	w.Flush()
	fmt.Printf("\nTotal: %d users\n", len(users))
	return nil
}

func Get(cfg *config.Config, username string, id int) error {
	token := cfg.Token

	if id != 0 {
		url := fmt.Sprintf("%s/api/v4/users/%d", cfg.GitLabURL, id)
		body, err := doRequest("GET", url, token, nil)
		if err != nil {
			return err
		}
		printJSON(body)
		return nil
	}
	if username != "" {
		url := fmt.Sprintf("%s/api/v4/users?username=%s", cfg.GitLabURL, username)
		body, err := doRequest("GET", url, token, nil)
		if err != nil {
			return err
		}
		var users []map[string]interface{}
		if err := json.Unmarshal(body, &users); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		if len(users) == 0 {
			fmt.Printf("User '%s' not found.\n", username)
			return nil
		}
		uid := intOrDash(users[0]["id"])
		url = fmt.Sprintf("%s/api/v4/users/%s", cfg.GitLabURL, uid)
		body, err = doRequest("GET", url, token, nil)
		if err != nil {
			return err
		}
		printJSON(body)
		return nil
	}
	return fmt.Errorf("--username or --id is required")
}

func Create(cfg *config.Config, username, email, name, password string, admin, skipConfirm, resetPw bool) error {
	token := cfg.Token

	payload := map[string]interface{}{
		"username":          username,
		"email":             email,
		"name":              name,
		"password":          password,
		"admin":             admin,
		"skip_confirmation": skipConfirm,
		"reset_password":    resetPw,
	}

	body, err := doRequest("POST", cfg.GitLabURL+"/api/v4/users", token, payload)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func SetPassword(cfg *config.Config, username string, id int, password string) error {
	token := cfg.Token

	userID := id
	if userID == 0 {
		if username == "" {
			return fmt.Errorf("--username or --id is required")
		}
		var err error
		userID, err = resolveUsername(cfg.GitLabURL, token, username)
		if err != nil {
			return err
		}
	}

	payload := map[string]interface{}{
		"password": password,
	}
	url := fmt.Sprintf("%s/api/v4/users/%d", cfg.GitLabURL, userID)
	body, err := doRequest("PUT", url, token, payload)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if msg, ok := result["message"]; ok {
			return fmt.Errorf("GitLab error: %v", msg)
		}
		fmt.Printf("Password updated for user '%v' (id=%d)\n", result["username"], userID)
	} else {
		fmt.Printf("Password updated for user id=%d\n", userID)
	}
	return nil
}

func Block(cfg *config.Config, username string, id int) error {
	token := cfg.Token

	userID, err := resolveUser(cfg.GitLabURL, token, id, username)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v4/users/%d/block", cfg.GitLabURL, userID)
	_, err = doRequest("POST", url, token, nil)
	if err != nil {
		return err
	}
	fmt.Printf("User id=%d blocked.\n", userID)
	return nil
}

func Unblock(cfg *config.Config, username string, id int) error {
	token := cfg.Token

	userID, err := resolveUser(cfg.GitLabURL, token, id, username)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v4/users/%d/unblock", cfg.GitLabURL, userID)
	_, err = doRequest("POST", url, token, nil)
	if err != nil {
		return err
	}
	fmt.Printf("User id=%d unblocked.\n", userID)
	return nil
}

func Delete(cfg *config.Config, username string, id int, hardDelete bool) error {
	token := cfg.Token

	userID, err := resolveUser(cfg.GitLabURL, token, id, username)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v4/users/%d", cfg.GitLabURL, userID)
	if hardDelete {
		url += "?hard_delete=true"
	}
	_, err = doRequest("DELETE", url, token, nil)
	if err != nil {
		return err
	}
	fmt.Printf("User id=%d deleted.\n", userID)
	return nil
}

func resolveUser(baseURL, token string, id int, username string) (int, error) {
	if id != 0 {
		return id, nil
	}
	if username != "" {
		return resolveUsername(baseURL, token, username)
	}
	return 0, fmt.Errorf("--username or --id is required")
}
