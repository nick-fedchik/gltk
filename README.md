# gltk

A fast, comprehensive GitLab CLI client built with Go and Cobra. Manage issues, merge requests, pipelines, CI/CD jobs, branches, tags, users, and more — all from the terminal.

## Install

```bash
# From source
go install github.com/gltk/gltk/cmd/gl@latest

# Or clone and build
git clone https://github.com/gltk/gltk.git
cd gltk
make build        # produces ./gl binary
make install      # installs to $GOPATH/bin/gl
```

**Requirements:** Go 1.25+

## Quick Start

```bash
# 1. Configure (pick one method)
export GITLAB_TOKEN="glpat-xxxxxxxxxxxxxxxxxxxx"

# Or create a config file
cp config.example.yaml config.yaml
# Edit config.yaml — set your token source

# 2. Verify connection
gl auth

# 3. Use it
gl issue list --project=mygroup/myproject
gl mr list --project=mygroup/myproject --state=opened
gl pipeline list --project=mygroup/myproject

# Or set a default project in config.yaml to skip --project on every command:
gl issue list
gl mr list --state=opened
gl pipeline list
```

## Configuration

The config format follows the official [gitlab.com/config/v1beta1](https://gitlab.com/gitlab-org/api/client-go/-/tree/main/config) schema — the same kubeconfig-style pattern used by the GitLab Go SDK.

### Config file search order

The `config.yaml` file is searched in this order (first found wins):

| Priority | Path | Use case |
|----------|------|----------|
| 1 | `./config.yaml` | Current working directory |
| 2 | `<bin-dir>/config.yaml` | Next to the `gl` binary |
| 3 | `.gltk/config.yaml` | Project-local config directory |
| 4 | `~/.config/gltk/config.yaml` | XDG user config (global) |

### Override precedence

| Priority | Source | Example |
|----------|--------|---------|
| 1 | CLI flags | `--token=glpat-...` `--gitlab-url=https://...` |
| 2 | Environment variables | `GITLAB_TOKEN`, `GITLAB_URL` |
| 3 | `config.yaml` | Resolved via search order above |
| 4 | Defaults | `https://gitlab.com` |

### Config file

```bash
cp config.example.yaml config.yaml
```

The config uses **instances** (servers), **auths** (credentials), and **contexts** (instance + auth pairs):

```yaml
version: gitlab.com/config/v1beta1
current-context: default

instances:
  - name: gitlab-com
    server: https://gitlab.com
    api-version: v4

auths:
  - name: my-token
    auth-info:
      personal-access-token:
        token-source:
          env-var: GITLAB_TOKEN     # read from $GITLAB_TOKEN

contexts:
  - name: default
    instance: gitlab-com
    auth: my-token
    default-project: "mygroup/myproject"  # optional — makes --project flag optional
```

### Authentication methods

**Personal Access Token** (recommended) — create at Settings > Access Tokens:

```yaml
auths:
  # From environment variable
  - name: token-env
    auth-info:
      personal-access-token:
        token-source:
          env-var: GITLAB_TOKEN

  # From file (newline is trimmed)
  - name: token-file
    auth-info:
      personal-access-token:
        token-source:
          file: ~/.config/gltk/token

  # Inline (only for local/private configs)
  - name: token-inline
    auth-info:
      personal-access-token:
        token: "glpat-xxxxxxxxxxxxxxxxxxxx"
```

**CI/CD Job Token** — for use in GitLab pipelines:

```yaml
auths:
  - name: ci
    auth-info:
      job-token:
        token-source:
          env-var: CI_JOB_TOKEN
```

**OAuth2** — for apps with OAuth flow:

```yaml
auths:
  - name: oauth
    auth-info:
      oauth2:
        access-token: "..."
        client-id: "your-app-id"
```

### Multiple instances

```yaml
version: gitlab.com/config/v1beta1
current-context: work

instances:
  - name: gitlab-com
    server: https://gitlab.com
  - name: company
    server: https://gitlab.company.com

auths:
  - name: personal
    auth-info:
      personal-access-token:
        token-source:
          env-var: GITLAB_TOKEN
  - name: work-token
    auth-info:
      personal-access-token:
        token-source:
          file: ~/.config/gltk/work-token

contexts:
  - name: personal
    instance: gitlab-com
    auth: personal
  - name: work
    instance: company
    auth: work-token
    default-project: "team/myapp"   # optional default project for this context
```

Switch context by editing `current-context`, or override per-command:

```bash
gl --gitlab-url=https://gitlab.company.com --token=glpat-... issue list --project=team/app
```

### Token scopes

Create a Personal Access Token at **Settings > Access Tokens** in your GitLab instance.

| Scope | Access Level |
|-------|-------------|
| `read_api` | Read-only (list issues, pipelines, MRs) |
| `api` | Full read/write (create/close issues, merge MRs) |
| `read_repository` | Read files and commits |
| `write_repository` | Push commits via API |
| `read_user` | Inspect users |
| `sudo` | Admin operations (manage users, runners) |

**Recommended minimum:** `api` covers most use cases.

## Commands

### Issues

```bash
gl issue list    --project=GROUP/PROJECT [--state=opened|closed]
gl issue create  --project=GROUP/PROJECT --title="Bug report" [--description="..."] [--labels="bug,urgent"]
                 [--assignee=USERNAME] [--assignee=USERNAME2]    # by username (resolved to ID automatically)
                 [--assignee-id=42] [--assignee-id=55]           # by numeric user ID
gl issue close   --project=GROUP/PROJECT --issue=42
gl issue batch   --project=GROUP/PROJECT --file=issues.json

# --project is optional on all commands when default-project is set in config.yaml
```

### Merge Requests

```bash
gl mr list       --project=GROUP/PROJECT [--state=opened|merged|closed]
gl mr get        --project=GROUP/PROJECT --mr=15
gl mr create     --project=GROUP/PROJECT --source=feature --target=main --title="Add feature"
                 [--assignee=USERNAME] [--assignee-id=42]        # assign by username or ID
gl mr merge      --project=GROUP/PROJECT --mr=15
gl mr close      --project=GROUP/PROJECT --mr=15
gl mr comment    --project=GROUP/PROJECT --mr=15 --body="LGTM"
```

### Pipelines & CI/CD

```bash
gl pipeline list        --project=GROUP/PROJECT
gl pipeline jobs        --project=GROUP/PROJECT --pipeline=123
gl pipeline watch       --project=GROUP/PROJECT --pipeline=123
gl pipeline cancel      --project=GROUP/PROJECT --pipeline=123
gl pipeline create      --project=GROUP/PROJECT --ref=main
gl pipeline trigger-job --project=GROUP/PROJECT --job=456
```

### Test Reports & Coverage

```bash
# Full test report with all test cases (JUnit)
gl pipeline test-report  --project=GROUP/PROJECT --pipeline=123

# Show only failed/errored tests
gl pipeline test-report  --project=GROUP/PROJECT --pipeline=123 --failed

# Quick summary — suite-level stats without individual tests
gl pipeline test-summary --project=GROUP/PROJECT --pipeline=123

# Coverage — pipeline total + per-job breakdown with visual bar
gl pipeline coverage     --project=GROUP/PROJECT --pipeline=123
```

Example output for `test-report --failed`:
```
Pipeline #123 — Test Report
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Total: 142  ✅ 139  ❌ 2  ⏭ 1  ⚠ 0
  Time:  34.50s

❌ api-tests (48 tests, 12.30s)
  ❌ TestUserAuth > should reject expired token (0.05s)
       expected status 401, got 200
       at auth_test.go:42
  ❌ TestUserAuth > should rate limit after 100 requests (1.20s)
       timeout after 5s
       ⚡ Failed 3 times on main
```

Example output for `coverage`:
```
Pipeline #123 — Coverage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Ref:    main
  Status: success
  Coverage: 74.20%

  Per-job coverage:
    unit-tests                      82.50%  [################----]
    integration-tests               61.30%  [############--------]
```

### Jobs

```bash
gl job analyze   --project=GROUP/PROJECT --job=789
gl job logs      --project=GROUP/PROJECT --job=789
gl job status    --project=GROUP/PROJECT --job=789
gl job details   --project=GROUP/PROJECT --job=789
gl job retry     --project=GROUP/PROJECT --job=789
gl job cancel    --project=GROUP/PROJECT --job=789
gl job trigger   --project=GROUP/PROJECT --job=789
gl job trace     --project=GROUP/PROJECT --job=789
```

### Branches

```bash
gl branch list     --project=GROUP/PROJECT
gl branch get      --project=GROUP/PROJECT --branch=main
gl branch create   --project=GROUP/PROJECT --branch=feature --ref=main
gl branch delete   --project=GROUP/PROJECT --branch=feature
gl branch protect  --project=GROUP/PROJECT --branch=main
```

### Tags

```bash
gl tag list      --project=GROUP/PROJECT
gl tag create    --project=GROUP/PROJECT --tag=v1.0.0 --ref=main [--message="Release"]
gl tag delete    --project=GROUP/PROJECT --tag=v1.0.0
```

### Files & Commits

```bash
gl file read     --project=GROUP/PROJECT --path=README.md [--ref=main]
gl file write    --project=GROUP/PROJECT --path=file.txt --content="..." --branch=main --message="Update"
gl file list     --project=GROUP/PROJECT [--path=src/] [--ref=main]
gl commit create --project=GROUP/PROJECT --file=commit-spec.json
gl diff summary  --project=GROUP/PROJECT --mr=15
```

### Comments

```bash
gl comment list    --project=GROUP/PROJECT --issue=42
gl comment add     --project=GROUP/PROJECT --issue=42 --body="Looks good"
gl comment delete  --project=GROUP/PROJECT --note=12345
```

### Labels & Milestones

```bash
gl label list       --project=GROUP/PROJECT
gl label create     --project=GROUP/PROJECT --name="priority::high" --color="#ff0000"
gl milestone list   --project=GROUP/PROJECT
gl milestone create --project=GROUP/PROJECT --title="v2.0"
gl milestone update --project=GROUP/PROJECT --milestone=1 --state=close
gl milestone delete --project=GROUP/PROJECT --milestone=1
```

### Users (admin)

```bash
gl user list
gl user get          --user=5
gl user create       --email=dev@example.com --username=dev --name="Developer"
gl user set-password --user=5 --password="newpass"
gl user block        --user=5
gl user unblock      --user=5
gl user delete       --user=5
```

### Runners (admin)

```bash
gl runner list     [--status=online|offline] [--type=instance|group|project]
gl runner get      --runner=10
gl runner delete   --runner=10
gl runner pause    --runner=10
gl runner resume   --runner=10
gl runner jobs     --runner=10
gl runner update   --runner=10 --description="Updated"
```

### Other

```bash
gl auth                  # Test authentication and connectivity
gl project list          # List accessible projects (membership)
gl project members       # List project members with access levels
gl search issues     --query="login bug" --project=GROUP/PROJECT
gl search mrs        --query="refactor" --project=GROUP/PROJECT
gl report            # Generate project reports
gl artifact list     --project=GROUP/PROJECT --job=789
gl artifact download --project=GROUP/PROJECT --job=789 [--output=artifacts.zip]
gl artifact delete   --project=GROUP/PROJECT --job=789
gl issues-check      --project=GROUP/PROJECT   # Interactive issue review
```

## Global Flags

All commands support these flags:

```
--gitlab-url    GitLab instance URL (overrides config)
--token         Personal access token (overrides config)
--format        Output format: text, json, table (default: text)
--verbose       Enable verbose output
--no-color      Disable color output
--config        Config file path (default: config.yaml)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GITLAB_URL` | GitLab instance URL |
| `GITLAB_TOKEN` | Personal access token (primary) |
| `GITLAB_API_TOKEN` | Personal access token (alternative) |
| `PRIVATE_TOKEN` | Personal access token (legacy) |

## Authors & Contributors

**Created by:** Nick Fedchyk (@nick-fedchik)

**Developed with assistance from:** Claude AI (Anthropic)
- Architecture design and feature planning
- Command implementation and testing
- Documentation and guides
- Code review and optimization

## License

Apache 2.0 — see [LICENSE](LICENSE).
