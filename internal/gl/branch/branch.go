package branch

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

func List(cfg *config.Config, project string, page, perPage int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListBranchesOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	branches, resp, err := client.Branches.ListBranches(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	if len(branches) == 0 {
		fmt.Println("No branches found")
		return nil
	}

	fmt.Printf("Branches (%d total):\n\n", resp.TotalItems)
	for _, b := range branches {
		protected := " "
		if b.Protected {
			protected = "🔒"
		}
		msg := b.Commit.Message
		if len(msg) > 50 {
			msg = msg[:50]
		}
		fmt.Printf("%s  %-40s  %s  %s\n",
			protected,
			b.Name,
			b.Commit.ShortID,
			msg,
		)
	}
	fmt.Printf("\nPage %d/%d\n", page, resp.TotalPages)
	return nil
}

func Get(cfg *config.Config, project, branchName string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	branchData, _, err := client.Branches.GetBranch(projectID, branchName)
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}

	fmt.Printf("Branch: %s\n", branchData.Name)
	fmt.Printf("Protected: %v\n", branchData.Protected)
	fmt.Printf("Commit: %s\n", branchData.Commit.ID)
	fmt.Printf("Message: %s\n", branchData.Commit.Message)
	fmt.Printf("Author: %s\n", branchData.Commit.AuthorName)
	fmt.Printf("Created: %s\n", branchData.Commit.CreatedAt.Format("2006-01-02 15:04:05"))
	return nil
}

func Create(cfg *config.Config, project, branchName, ref string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.CreateBranchOptions{
		Branch: glclient.Ptr(branchName),
		Ref:    glclient.Ptr(ref),
	}

	branchData, _, err := client.Branches.CreateBranch(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	fmt.Printf("✓ Branch created\n")
	fmt.Printf("  Name: %s\n", branchData.Name)
	fmt.Printf("  Commit: %s\n", branchData.Commit.ShortID)
	return nil
}

func Delete(cfg *config.Config, project, branchName string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	_, err = client.Branches.DeleteBranch(projectID, branchName)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	fmt.Printf("✓ Branch deleted\n")
	fmt.Printf("  Name: %s\n", branchName)
	return nil
}

func Protect(cfg *config.Config, project, branchName string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	branchData, _, err := client.Branches.GetBranch(projectID, branchName)
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}

	fmt.Printf("Branch Protection Status\n")
	fmt.Printf("  Name: %s\n", branchData.Name)
	if branchData.Protected {
		fmt.Printf("  Status: 🔒 Protected\n")
	} else {
		fmt.Printf("  Status: 🔓 Unprotected\n")
	}
	fmt.Printf("\nNote: To modify protection rules, use GitLab web interface:\n")
	fmt.Printf("  Settings → Repository → Protected branches\n")
	return nil
}
