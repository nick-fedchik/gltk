# GL Test Reports & Coverage — Guide

## Prerequisites

GitLab automatically parses JUnit XML artifacts from CI jobs. No Java required — any language can generate the format.

### CI Configuration

Add `artifacts:reports:junit` to your `.gitlab-ci.yml`:

```yaml
# Go
test:
  script:
    - gotestsum --junitfile report.xml ./...
  artifacts:
    reports:
      junit: report.xml

# Node.js (Jest)
test:
  script:
    - npx jest --reporters=jest-junit
  artifacts:
    reports:
      junit: junit.xml

# Node.js (Vitest)
test:
  script:
    - npx vitest --reporter=junit --outputFile=junit.xml
  artifacts:
    reports:
      junit: junit.xml

# Python
test:
  script:
    - pytest --junitxml=report.xml
  artifacts:
    reports:
      junit: report.xml
```

### Coverage Configuration

Coverage is extracted from job stdout via regex. Configure in **Settings > CI/CD > General pipelines > Test coverage parsing**, or in `.gitlab-ci.yml`:

```yaml
test:
  script:
    - go test -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  coverage: '/total:\s+\(statements\)\s+(\d+\.\d+)%/'
```

Common regex patterns:

| Language | Regex |
|----------|-------|
| Go | `/total:\s+\(statements\)\s+(\d+\.\d+)%/` |
| Jest | `/All files\s*\|\s*(\d+\.?\d*)/` |
| pytest-cov | `/TOTAL\s+\d+\s+\d+\s+(\d+)%/` |
| Istanbul | `/Statements\s*:\s*(\d+\.?\d*)%/` |

---

## Commands

### gl pipeline test-report

Full test report with individual test cases, stack traces, and failure history.

```bash
# All tests
gl pipeline test-report --project=GROUP/PROJECT --pipeline=123

# Only failed and errored tests
gl pipeline test-report --project=GROUP/PROJECT --pipeline=123 --failed
```

Output:
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

Fields per test case:
- **Status** — success, failed, error, skipped
- **Classname > Name** — test hierarchy
- **ExecutionTime** — duration in seconds
- **StackTrace** — first 10 lines of failure output
- **RecentFailures** — flaky test indicator (how many times it failed on default branch)

### gl pipeline test-summary

Lightweight suite-level summary without individual test cases. Fast — does not load all test case data.

```bash
gl pipeline test-summary --project=GROUP/PROJECT --pipeline=123
```

Output:
```
Pipeline #123 — Test Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Total: 142  ✅ 139  ❌ 2  ⏭ 1  ⚠ 0  (34.50s)

  ✅ unit-tests                               94 total   94 pass    0 fail  (22.20s)
  ❌ api-tests                                48 total   45 pass    2 fail  (12.30s)
```

### gl pipeline coverage

Pipeline-level and per-job coverage breakdown with ASCII visualization.

```bash
gl pipeline coverage --project=GROUP/PROJECT --pipeline=123
```

Output:
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

---

## Creating Issues with Assignees

Create issues and assign them to specific users:

```bash
# Get user IDs first
gl user list --project=GROUP/PROJECT

# Create and assign to single user
gl issue create --project=GROUP/PROJECT --title="Fix failing test" --assignee=42

# Create and assign to multiple users (GitLab EE)
gl issue create --project=GROUP/PROJECT --title="Fix failing test" --assignee=42 --assignee=55

# With description and labels
gl issue create --project=GROUP/PROJECT \
  --title="Fix failing test" \
  --description="Failing test: TestUserAuth > should reject expired token" \
  --labels="bug,test-failure" \
  --assignee=42
```

Output:
```
✅ Created issue #456: Fix failing test
   Assigned to: john_doe, jane_smith
```

---

## API Reference

| Command | GitLab API | Response size |
|---------|-----------|---------------|
| `test-report` | `GET /projects/:id/pipelines/:pid/test_report` | Full — all test cases with details |
| `test-summary` | `GET /projects/:id/pipelines/:pid/test_report_summary` | Light — suite-level aggregates only |
| `coverage` | `GET /projects/:id/pipelines/:pid` + `GET /projects/:id/pipelines/:pid/jobs` | Pipeline metadata + job list |

## Typical Workflow

```bash
# 1. Check what pipelines ran
gl pipeline list --project=mygroup/myapp

# 2. Quick check — did tests pass?
gl pipeline test-summary --project=mygroup/myapp --pipeline=456

# 3. Something failed — see details
gl pipeline test-report --project=mygroup/myapp --pipeline=456 --failed

# 4. Check coverage
gl pipeline coverage --project=mygroup/myapp --pipeline=456

# 5. Look at the failing job's full log
gl pipeline trace --project=mygroup/myapp --job=789
```
