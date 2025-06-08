# ESH CLI Integrations & Extensions

## 🔗 Integration Opportunities

### 1. GitHub/GitLab Integration

#### GitHub Releases Integration
```go
// pkg/integrations/github.go
package integrations

import (
    "context"
    "github.com/google/go-github/v50/github"
    "golang.org/x/oauth2"
)

type GitHubIntegration struct {
    client *github.Client
    owner  string
    repo   string
}

func NewGitHubIntegration(token, owner, repo string) *GitHubIntegration {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(context.Background(), ts)
    
    return &GitHubIntegration{
        client: github.NewClient(tc),
        owner:  owner,
        repo:   repo,
    }
}

func (g *GitHubIntegration) CreateRelease(tag, changelog string) error {
    release := &github.RepositoryRelease{
        TagName:    &tag,
        Name:       &tag,
        Body:       &changelog,
        Draft:      github.Bool(false),
        Prerelease: github.Bool(false),
    }
    
    _, _, err := g.client.Repositories.CreateRelease(
        context.Background(), g.owner, g.repo, release,
    )
    return err
}

func (g *GitHubIntegration) GetLatestRelease() (*github.RepositoryRelease, error) {
    release, _, err := g.client.Repositories.GetLatestRelease(
        context.Background(), g.owner, g.repo,
    )
    return release, err
}
```

#### GitHub Actions Integration
```yaml
# .github/workflows/auto-release.yml
name: Auto Release

on:
  push:
    tags:
      - '*_production2_*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - name: Install ESH CLI
      run: |
        curl -L https://github.com/PocketfulDev/esh-cli/releases/latest/download/esh-cli-linux-amd64 -o esh-cli
        chmod +x esh-cli
    
    - name: Generate Changelog
      id: changelog
      run: |
        CHANGELOG=$(./esh-cli changelog --conventional-commits --format markdown)
        echo "changelog<<EOF" >> $GITHUB_OUTPUT
        echo "$CHANGELOG" >> $GITHUB_OUTPUT
        echo "EOF" >> $GITHUB_OUTPUT
    
    - name: Create GitHub Release
      uses: actions/create-release@v1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with:
        tag_name: ${{ github.ref_name }}
        release_name: ${{ github.ref_name }}
        body: ${{ steps.changelog.outputs.changelog }}
        draft: false
        prerelease: false
```

### 2. Slack/Teams Integration

#### Slack Notifications
```go
// pkg/integrations/slack.go
package integrations

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type SlackIntegration struct {
    webhookURL string
}

type SlackMessage struct {
    Text        string            `json:"text"`
    Username    string            `json:"username"`
    IconEmoji   string            `json:"icon_emoji"`
    Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
    Color  string       `json:"color"`
    Fields []SlackField `json:"fields"`
}

type SlackField struct {
    Title string `json:"title"`
    Value string `json:"value"`
    Short bool   `json:"short"`
}

func NewSlackIntegration(webhookURL string) *SlackIntegration {
    return &SlackIntegration{webhookURL: webhookURL}
}

func (s *SlackIntegration) NotifyVersionBump(tag, environment, bumpType string) error {
    message := SlackMessage{
        Text:      fmt.Sprintf("🚀 New version deployed: %s", tag),
        Username:  "ESH CLI",
        IconEmoji: ":rocket:",
        Attachments: []SlackAttachment{
            {
                Color: getColorForBumpType(bumpType),
                Fields: []SlackField{
                    {Title: "Environment", Value: environment, Short: true},
                    {Title: "Version", Value: tag, Short: true},
                    {Title: "Bump Type", Value: bumpType, Short: true},
                    {Title: "Timestamp", Value: time.Now().Format("2006-01-02 15:04:05"), Short: true},
                },
            },
        },
    }
    
    return s.sendMessage(message)
}

func (s *SlackIntegration) sendMessage(message SlackMessage) error {
    jsonData, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("slack API returned status %d", resp.StatusCode)
    }
    
    return nil
}

func getColorForBumpType(bumpType string) string {
    switch bumpType {
    case "MAJOR":
        return "danger"  // Red
    case "MINOR":
        return "good"    // Green
    case "PATCH":
        return "warning" // Yellow
    default:
        return "#36a64f" // Default green
    }
}
```

#### Teams Integration
```go
// pkg/integrations/teams.go
package integrations

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type TeamsIntegration struct {
    webhookURL string
}

type TeamsMessage struct {
    Type       string                 `json:"@type"`
    Context    string                 `json:"@context"`
    Summary    string                 `json:"summary"`
    Title      string                 `json:"title"`
    Text       string                 `json:"text"`
    ThemeColor string                 `json:"themeColor"`
    Sections   []TeamsSection         `json:"sections,omitempty"`
    Actions    []TeamsAction          `json:"potentialAction,omitempty"`
}

type TeamsSection struct {
    ActivityTitle    string      `json:"activityTitle"`
    ActivitySubtitle string      `json:"activitySubtitle"`
    Facts           []TeamsFact `json:"facts"`
}

type TeamsFact struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type TeamsAction struct {
    Type    string `json:"@type"`
    Name    string `json:"name"`
    Targets []struct {
        OS  string `json:"os"`
        URI string `json:"uri"`
    } `json:"targets"`
}

func NewTeamsIntegration(webhookURL string) *TeamsIntegration {
    return &TeamsIntegration{webhookURL: webhookURL}
}

func (t *TeamsIntegration) NotifyVersionBump(tag, environment, bumpType, changelog string) error {
    message := TeamsMessage{
        Type:       "MessageCard",
        Context:    "http://schema.org/extensions",
        Summary:    fmt.Sprintf("Version %s deployed to %s", tag, environment),
        Title:      "🚀 ESH CLI Version Deployment",
        Text:       fmt.Sprintf("New %s version deployed", bumpType),
        ThemeColor: getTeamsColorForBumpType(bumpType),
        Sections: []TeamsSection{
            {
                ActivityTitle:    "Version Details",
                ActivitySubtitle: fmt.Sprintf("Deployed to %s environment", environment),
                Facts: []TeamsFact{
                    {Name: "Tag", Value: tag},
                    {Name: "Environment", Value: environment},
                    {Name: "Bump Type", Value: bumpType},
                    {Name: "Deployed At", Value: time.Now().Format("2006-01-02 15:04:05 MST")},
                },
            },
        },
    }
    
    if changelog != "" {
        message.Sections = append(message.Sections, TeamsSection{
            ActivityTitle: "Changelog",
            Facts: []TeamsFact{
                {Name: "Changes", Value: changelog},
            },
        })
    }
    
    return t.sendMessage(message)
}
```

### 3. JIRA Integration

#### JIRA Ticket Linking
```go
// pkg/integrations/jira.go
package integrations

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "strings"
)

type JIRAIntegration struct {
    baseURL  string
    username string
    token    string
    client   *http.Client
}

type JIRAIssue struct {
    ID     string            `json:"id"`
    Key    string            `json:"key"`
    Fields JIRAIssueFields   `json:"fields"`
}

type JIRAIssueFields struct {
    Summary     string        `json:"summary"`
    Description string        `json:"description"`
    Status      JIRAStatus    `json:"status"`
    FixVersions []JIRAVersion `json:"fixVersions,omitempty"`
}

type JIRAStatus struct {
    Name string `json:"name"`
}

type JIRAVersion struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Released    bool   `json:"released"`
}

func NewJIRAIntegration(baseURL, username, token string) *JIRAIntegration {
    return &JIRAIntegration{
        baseURL:  strings.TrimSuffix(baseURL, "/"),
        username: username,
        token:    token,
        client:   &http.Client{},
    }
}

func (j *JIRAIntegration) ExtractTicketNumbers(commits []string) []string {
    ticketPattern := regexp.MustCompile(`([A-Z]+-\d+)`)
    ticketSet := make(map[string]bool)
    
    for _, commit := range commits {
        matches := ticketPattern.FindAllString(commit, -1)
        for _, match := range matches {
            ticketSet[match] = true
        }
    }
    
    tickets := make([]string, 0, len(ticketSet))
    for ticket := range ticketSet {
        tickets = append(tickets, ticket)
    }
    return tickets
}

func (j *JIRAIntegration) AddFixVersion(ticketKey, version string) error {
    issue, err := j.getIssue(ticketKey)
    if err != nil {
        return err
    }
    
    // Add new fix version
    newVersion := JIRAVersion{
        Name:     version,
        Released: true,
    }
    
    issue.Fields.FixVersions = append(issue.Fields.FixVersions, newVersion)
    
    return j.updateIssue(ticketKey, issue)
}

func (j *JIRAIntegration) AddComment(ticketKey, comment string) error {
    commentData := map[string]interface{}{
        "body": comment,
    }
    
    return j.makeRequest("POST", 
        fmt.Sprintf("/rest/api/2/issue/%s/comment", ticketKey), 
        commentData, nil)
}

func (j *JIRAIntegration) getIssue(key string) (*JIRAIssue, error) {
    var issue JIRAIssue
    err := j.makeRequest("GET", 
        fmt.Sprintf("/rest/api/2/issue/%s", key), 
        nil, &issue)
    return &issue, err
}

func (j *JIRAIntegration) makeRequest(method, endpoint string, body interface{}, result interface{}) error {
    // Implementation details for JIRA API requests
    // Including authentication, request/response handling
    return nil
}
```

### 4. CI/CD Pipeline Integration

#### Jenkins Integration
```groovy
// Jenkinsfile
pipeline {
    agent any
    
    environment {
        ESH_CLI_VERSION = 'latest'
    }
    
    stages {
        stage('Install ESH CLI') {
            steps {
                sh 'curl -L https://github.com/PocketfulDev/esh-cli/releases/latest/download/esh-cli-linux-amd64 -o esh-cli'
                sh 'chmod +x esh-cli'
            }
        }
        
        stage('Version Analysis') {
            steps {
                script {
                    // Get branch-based version suggestions
                    def suggestion = sh(
                        script: './esh-cli branch-version --suggest',
                        returnStdout: true
                    ).trim()
                    
                    // Store suggestion for later stages
                    env.VERSION_SUGGESTION = suggestion
                }
            }
        }
        
        stage('Deploy to Dev') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'develop') {
                        sh './esh-cli branch-version --auto-tag --env dev'
                    }
                }
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'release/*'
            }
            steps {
                sh './esh-cli bump-version stg6 --auto'
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                script {
                    // Get latest staging tag for promotion
                    def latestStg6 = sh(
                        script: './esh-cli last-tag stg6',
                        returnStdout: true
                    ).trim()
                    
                    // Promote to production
                    sh "./esh-cli add-tag production2 ${latestStg6} --from stg6_${latestStg6}"
                    
                    // Generate release notes
                    def changelog = sh(
                        script: './esh-cli changelog production2 --conventional-commits --format markdown',
                        returnStdout: true
                    ).trim()
                    
                    // Archive changelog
                    writeFile file: 'RELEASE_NOTES.md', text: changelog
                    archiveArtifacts artifacts: 'RELEASE_NOTES.md'
                }
            }
        }
    }
    
    post {
        success {
            script {
                if (env.BRANCH_NAME == 'main') {
                    // Notify Slack of production deployment
                    slackSend(
                        channel: '#deployments',
                        color: 'good',
                        message: "🚀 Production deployment successful: ${env.BUILD_URL}"
                    )
                }
            }
        }
        
        failure {
            slackSend(
                channel: '#deployments',
                color: 'danger',
                message: "❌ Deployment failed: ${env.BUILD_URL}"
            )
        }
    }
}
```

#### GitLab CI Integration
```yaml
# .gitlab-ci.yml
stages:
  - analyze
  - deploy-dev
  - deploy-staging
  - deploy-production

variables:
  ESH_CLI_URL: "https://github.com/PocketfulDev/esh-cli/releases/latest/download/esh-cli-linux-amd64"

before_script:
  - curl -L $ESH_CLI_URL -o esh-cli
  - chmod +x esh-cli

analyze-version:
  stage: analyze
  script:
    - ./esh-cli branch-version --suggest
    - ./esh-cli version-diff dev --history
  artifacts:
    reports:
      dotenv: version-analysis.env

deploy-dev:
  stage: deploy-dev
  script:
    - ./esh-cli branch-version --auto-tag --env dev
  only:
    - develop
  environment:
    name: development

deploy-staging:
  stage: deploy-staging
  script:
    - ./esh-cli bump-version stg6 --auto
    - ./esh-cli changelog stg6 --conventional-commits --output changelog.md
  artifacts:
    paths:
      - changelog.md
  only:
    - /^release\/.*$/
  environment:
    name: staging

deploy-production:
  stage: deploy-production
  script:
    - LATEST_STG6=$(./esh-cli last-tag stg6)
    - ./esh-cli add-tag production2 $LATEST_STG6 --from stg6_$LATEST_STG6
    - ./esh-cli changelog production2 --full --format markdown --output RELEASE_NOTES.md
  artifacts:
    paths:
      - RELEASE_NOTES.md
  only:
    - main
  environment:
    name: production
```

### 5. Monitoring and Analytics Integration

#### Datadog Integration
```go
// pkg/integrations/datadog.go
package integrations

import (
    "github.com/DataDog/datadog-go/statsd"
)

type DatadogIntegration struct {
    client *statsd.Client
}

func NewDatadogIntegration(addr string) (*DatadogIntegration, error) {
    client, err := statsd.New(addr)
    if err != nil {
        return nil, err
    }
    
    return &DatadogIntegration{client: client}, nil
}

func (d *DatadogIntegration) RecordVersionBump(environment, bumpType string) error {
    tags := []string{
        fmt.Sprintf("environment:%s", environment),
        fmt.Sprintf("bump_type:%s", bumpType),
    }
    
    // Increment deployment counter
    err := d.client.Incr("esh_cli.version_bump", tags, 1)
    if err != nil {
        return err
    }
    
    // Record deployment timing
    return d.client.Timing("esh_cli.deployment_time", 
        time.Since(startTime), tags, 1)
}

func (d *DatadogIntegration) RecordVersionMetrics(environment string, 
    majorCount, minorCount, patchCount int) error {
    
    tags := []string{fmt.Sprintf("environment:%s", environment)}
    
    d.client.Gauge("esh_cli.major_versions", float64(majorCount), tags, 1)
    d.client.Gauge("esh_cli.minor_versions", float64(minorCount), tags, 1)
    d.client.Gauge("esh_cli.patch_versions", float64(patchCount), tags, 1)
    
    return nil
}
```

#### Prometheus Integration
```go
// pkg/integrations/prometheus.go
package integrations

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    versionBumps = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "esh_cli_version_bumps_total",
            Help: "Total number of version bumps by environment and type",
        },
        []string{"environment", "bump_type"},
    )
    
    deploymentDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "esh_cli_deployment_duration_seconds",
            Help: "Duration of deployment operations",
        },
        []string{"environment", "operation"},
    )
    
    activeVersions = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "esh_cli_active_versions",
            Help: "Number of active versions per environment",
        },
        []string{"environment"},
    )
)

type PrometheusIntegration struct{}

func NewPrometheusIntegration() *PrometheusIntegration {
    return &PrometheusIntegration{}
}

func (p *PrometheusIntegration) RecordVersionBump(environment, bumpType string) {
    versionBumps.WithLabelValues(environment, bumpType).Inc()
}

func (p *PrometheusIntegration) RecordDeploymentDuration(environment, operation string, duration time.Duration) {
    deploymentDuration.WithLabelValues(environment, operation).Observe(duration.Seconds())
}

func (p *PrometheusIntegration) UpdateActiveVersionCount(environment string, count float64) {
    activeVersions.WithLabelValues(environment).Set(count)
}
```

### 6. Configuration Management

#### Integration Configuration
```yaml
# ~/.esh-cli.yaml
integrations:
  slack:
    enabled: true
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#deployments"
    notify_on:
      - version_bump
      - deployment
      - error
  
  teams:
    enabled: false
    webhook_url: "${TEAMS_WEBHOOK_URL}"
  
  jira:
    enabled: true
    base_url: "https://company.atlassian.net"
    username: "${JIRA_USERNAME}"
    token: "${JIRA_API_TOKEN}"
    auto_link_tickets: true
    auto_add_fix_versions: true
  
  github:
    enabled: true
    token: "${GITHUB_TOKEN}"
    owner: "PocketfulDev"
    repo: "esh-cli"
    auto_create_releases: true
  
  monitoring:
    datadog:
      enabled: true
      agent_address: "localhost:8125"
    
    prometheus:
      enabled: true
      port: 9090

notifications:
  version_bump:
    - slack
    - teams
  
  deployment:
    - slack
    - jira
  
  error:
    - slack
    - monitoring
```

### 7. Integration Commands

Create integration management commands:

```go
// cmd/integrations.go
package cmd

import (
    "esh-cli/pkg/integrations"
    "github.com/spf13/cobra"
)

var integrationsCmd = &cobra.Command{
    Use:   "integrations",
    Short: "Manage external integrations",
    Long:  "Configure and test external integrations like Slack, JIRA, GitHub, etc.",
}

var testIntegrationsCmd = &cobra.Command{
    Use:   "test [integration]",
    Short: "Test integration connectivity",
    Args:  cobra.MaximumNArgs(1),
    Run:   runTestIntegrations,
}

var listIntegrationsCmd = &cobra.Command{
    Use:   "list",
    Short: "List available integrations",
    Run:   runListIntegrations,
}

func init() {
    rootCmd.AddCommand(integrationsCmd)
    integrationsCmd.AddCommand(testIntegrationsCmd)
    integrationsCmd.AddCommand(listIntegrationsCmd)
}

func runTestIntegrations(cmd *cobra.Command, args []string) {
    // Test integration connectivity
    if len(args) == 0 {
        // Test all enabled integrations
        testAllIntegrations()
    } else {
        // Test specific integration
        testSpecificIntegration(args[0])
    }
}

func runListIntegrations(cmd *cobra.Command, args []string) {
    // List all available integrations and their status
    fmt.Println("Available Integrations:")
    fmt.Println("🔗 Slack: " + getIntegrationStatus("slack"))
    fmt.Println("🔗 Teams: " + getIntegrationStatus("teams"))
    fmt.Println("🔗 JIRA: " + getIntegrationStatus("jira"))
    fmt.Println("🔗 GitHub: " + getIntegrationStatus("github"))
    fmt.Println("📊 Datadog: " + getIntegrationStatus("datadog"))
    fmt.Println("📊 Prometheus: " + getIntegrationStatus("prometheus"))
}
```

This comprehensive integration system transforms your ESH CLI into an enterprise-ready tool that seamlessly connects with your entire development and deployment ecosystem!
