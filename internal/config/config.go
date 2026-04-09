package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/gitlab-org/api/client-go"
	"gopkg.in/yaml.v3"
)

// Config is the resolved runtime configuration used by commands.
type Config struct {
	GitLabURL string
	Token     string
	ProjectID string
}

// ConfigFile represents the config.yaml on-disk format (kubeconfig-style).
type ConfigFile struct {
	Version        string       `yaml:"version"`
	CurrentContext string       `yaml:"current-context"`
	Preferences    Preferences  `yaml:"preferences,omitempty"`
	Instances      []Instance   `yaml:"instances"`
	Auths          []Auth       `yaml:"auths"`
	Contexts       []Context    `yaml:"contexts"`
}

// Preferences contains global client preferences.
type Preferences struct {
	RetryMax     int    `yaml:"retry-max,omitempty"`
	RetryWaitMin string `yaml:"retry-wait-min,omitempty"`
	RetryWaitMax string `yaml:"retry-wait-max,omitempty"`
}

// Instance represents a GitLab server.
type Instance struct {
	Name                  string    `yaml:"name"`
	Server                string    `yaml:"server"`
	APIVersion            string    `yaml:"api-version,omitempty"`
	InsecureSkipTLSVerify bool      `yaml:"insecure-skip-tls-verify,omitempty"`
	RateLimit             *RateLimit `yaml:"rate-limit,omitempty"`
}

// RateLimit contains rate limiting configuration.
type RateLimit struct {
	RequestsPerSecond float64 `yaml:"requests-per-second,omitempty"`
	Burst             int     `yaml:"burst,omitempty"`
}

// Auth contains authentication information.
type Auth struct {
	Name     string   `yaml:"name"`
	AuthInfo AuthInfo `yaml:"auth-info"`
}

// AuthInfo contains the authentication method details.
type AuthInfo struct {
	PersonalAccessToken *PersonalAccessToken `yaml:"personal-access-token,omitempty"`
	JobToken            *TokenField          `yaml:"job-token,omitempty"`
	OAuth2              *OAuth2              `yaml:"oauth2,omitempty"`
}

// PersonalAccessToken provides a PAT either directly or via a source.
type PersonalAccessToken struct {
	Token       string       `yaml:"token,omitempty"`
	TokenSource *TokenSource `yaml:"token-source,omitempty"`
}

// TokenField provides a generic token either directly or via a source.
type TokenField struct {
	Token       string       `yaml:"token,omitempty"`
	TokenSource *TokenSource `yaml:"token-source,omitempty"`
}

// OAuth2 contains OAuth2 authentication settings.
type OAuth2 struct {
	AccessToken  string       `yaml:"access-token,omitempty"`
	RefreshToken string       `yaml:"refresh-token,omitempty"`
	ClientID     string       `yaml:"client-id,omitempty"`
	ClientSecret string       `yaml:"client-secret,omitempty"`
	ClientSecretSource *TokenSource `yaml:"client-secret-source,omitempty"`
}

// TokenSource specifies how to obtain a credential value.
type TokenSource struct {
	EnvVar  string          `yaml:"env-var,omitempty"`
	File    string          `yaml:"file,omitempty"`
	Keyring *KeyringSource  `yaml:"keyring,omitempty"`
}

// KeyringSource specifies keyring/keychain credential storage.
type KeyringSource struct {
	Service string `yaml:"service"`
	User    string `yaml:"user"`
}

// Context binds an instance to an auth.
type Context struct {
	Name           string `yaml:"name"`
	Instance       string `yaml:"instance"`
	Auth           string `yaml:"auth"`
	DefaultProject string `yaml:"default-project,omitempty"`
}

// Load resolves configuration from env vars and config file.
func Load() (*Config, error) {
	cfg := &Config{
		GitLabURL: os.Getenv("GITLAB_URL"),
		Token:     firstNonEmpty(os.Getenv("GITLAB_TOKEN"), os.Getenv("GITLAB_API_TOKEN"), os.Getenv("PRIVATE_TOKEN")),
	}

	// Load from config file and resolve the active context
	if cf := tryLoadConfigFile(); cf != nil {
		resolved := cf.resolve()
		if cfg.GitLabURL == "" {
			cfg.GitLabURL = resolved.GitLabURL
		}
		if cfg.Token == "" {
			cfg.Token = resolved.Token
		}
		if cfg.ProjectID == "" {
			cfg.ProjectID = resolved.ProjectID
		}
	}

	if cfg.GitLabURL == "" {
		cfg.GitLabURL = "https://gitlab.com"
	}

	return cfg, nil
}

// resolve picks the active context and returns a flat Config.
func (cf *ConfigFile) resolve() *Config {
	ctxName := cf.CurrentContext
	if ctxName == "" && len(cf.Contexts) > 0 {
		ctxName = cf.Contexts[0].Name
	}

	var ctx *Context
	for i := range cf.Contexts {
		if cf.Contexts[i].Name == ctxName {
			ctx = &cf.Contexts[i]
			break
		}
	}
	if ctx == nil {
		return &Config{}
	}

	var inst *Instance
	for i := range cf.Instances {
		if cf.Instances[i].Name == ctx.Instance {
			inst = &cf.Instances[i]
			break
		}
	}

	var auth *Auth
	for i := range cf.Auths {
		if cf.Auths[i].Name == ctx.Auth {
			auth = &cf.Auths[i]
			break
		}
	}

	cfg := &Config{}
	if inst != nil {
		cfg.GitLabURL = inst.Server
	}
	if auth != nil {
		cfg.Token = resolveToken(auth)
	}
	if ctx.DefaultProject != "" {
		cfg.ProjectID = ctx.DefaultProject
	}
	return cfg
}

// resolveToken extracts a token string from an Auth entry.
func resolveToken(auth *Auth) string {
	ai := &auth.AuthInfo

	if ai.PersonalAccessToken != nil {
		if ai.PersonalAccessToken.Token != "" {
			return ai.PersonalAccessToken.Token
		}
		if ai.PersonalAccessToken.TokenSource != nil {
			return resolveSource(ai.PersonalAccessToken.TokenSource)
		}
	}

	if ai.JobToken != nil {
		if ai.JobToken.Token != "" {
			return ai.JobToken.Token
		}
		if ai.JobToken.TokenSource != nil {
			return resolveSource(ai.JobToken.TokenSource)
		}
	}

	if ai.OAuth2 != nil && ai.OAuth2.AccessToken != "" {
		return ai.OAuth2.AccessToken
	}

	return ""
}

// resolveSource reads a credential from an env var or file.
func resolveSource(src *TokenSource) string {
	if src.EnvVar != "" {
		return os.Getenv(src.EnvVar)
	}
	if src.File != "" {
		data, err := os.ReadFile(src.File)
		if err == nil {
			return trimNewline(string(data))
		}
	}
	// Keyring support requires OS-specific integration — not resolved here.
	return ""
}

func trimNewline(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// tryLoadConfigFile searches for config.yaml in this order:
//  1. ./config.yaml          — current working directory
//  2. <bin-dir>/config.yaml  — next to the gl binary
//  3. .gltk/config.yaml      — project-local config directory
//  4. ~/.config/gltk/config.yaml — XDG user config
func tryLoadConfigFile() *ConfigFile {
	const filename = "config.yaml"

	// 1. Current working directory
	if dir, err := os.Getwd(); err == nil {
		if cf := tryParse(filepath.Join(dir, filename)); cf != nil {
			return cf
		}
	}

	// 2. Next to the binary (bin/)
	if exe, err := os.Executable(); err == nil {
		binDir := filepath.Dir(exe)
		if cf := tryParse(filepath.Join(binDir, filename)); cf != nil {
			return cf
		}
	}

	// 3. .gltk/ in current directory
	if dir, err := os.Getwd(); err == nil {
		if cf := tryParse(filepath.Join(dir, ".gltk", filename)); cf != nil {
			return cf
		}
	}

	// 4. XDG config: ~/.config/gltk/config.yaml
	if home, err := os.UserHomeDir(); err == nil {
		if cf := tryParse(filepath.Join(home, ".config", "gltk", filename)); cf != nil {
			return cf
		}
	}

	return nil
}

func tryParse(path string) *ConfigFile {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cf ConfigFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: found %s but failed to parse: %v\n", path, err)
		return nil
	}
	return &cf
}

// NewGitLabClient creates a new GitLab API client from resolved config.
func (c *Config) NewGitLabClient() (*gitlab.Client, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("GitLab token is required. Set GITLAB_TOKEN env var or configure config.yaml")
	}
	client, err := gitlab.NewClient(c.Token, gitlab.WithBaseURL(c.GitLabURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return client, nil
}
