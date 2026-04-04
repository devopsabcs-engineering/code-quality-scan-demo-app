# Workshop Setup Guide

This guide covers prerequisites and setup for the Code Quality Scan Workshop.

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Git | 2.40+ | Source control |
| GitHub CLI (`gh`) | 2.40+ | GitHub API interaction |
| Azure CLI (`az`) | 2.55+ | Azure resource management |
| Docker Desktop | 24+ | Container builds and local testing |
| Node.js | 20 LTS | TypeScript demo app |
| Python | 3.12+ | Python demo app and SARIF converters |
| .NET SDK | 8.0+ | C# demo app |
| Java JDK | 21+ | Java demo app |
| Go | 1.22+ | Go demo app |
| PowerShell Core | 7.4+ | Script execution |

## Azure Requirements

1. **Azure Subscription** — Active subscription with Contributor access.
2. **Azure AD** — Permission to create app registrations.
3. **Azure Container Registry** — Will be provisioned by Bicep templates.
4. **Azure App Service** — Web App for Containers (Linux, B1 SKU).

## GitHub Requirements

1. **GitHub Account** — With access to the `devopsabcs-engineering` organization.
2. **GitHub Codespaces** — Recommended for consistent development environment.
3. **GitHub Advanced Security** — Required for code scanning features.

## Setup Steps

### 1. Fork/Clone the Scanner Repo

```powershell
gh repo fork devopsabcs-engineering/code-quality-scan-demo-app --clone
Set-Location code-quality-scan-demo-app
```

### 2. Run OIDC Setup

```powershell
.\scripts\setup-oidc.ps1 -SubscriptionId "<your-subscription-id>"
```

This creates:
- Azure AD app registration (`code-quality-scan-demo-app-oidc`)
- Federated identity credentials for GitHub Actions
- Service principal with Contributor role

### 3. Bootstrap Demo App Repos

```powershell
.\scripts\bootstrap-demo-apps.ps1 `
    -ClientId "<client-id-from-step-2>" `
    -TenantId "<tenant-id>" `
    -SubscriptionId "<subscription-id>"
```

This creates:
- 5 GitHub repositories (cq-demo-app-001 through 005)
- OIDC secrets on all repos
- `prod` environment on all repos

### 4. Deploy Demo Apps

Either run the `deploy-all.yml` workflow from the scanner repo, or deploy individually:

```powershell
# Trigger from GitHub UI or CLI
gh workflow run deploy-all.yml --repo devopsabcs-engineering/code-quality-scan-demo-app
```

### 5. Verify Deployment

```powershell
# Check each app's health endpoint
$apps = @("001","002","003","004","005")
foreach ($app in $apps) {
    $rg = "rg-cq-demo-$app"
    $name = az webapp list -g $rg --query "[0].defaultHostName" -o tsv
    Write-Host "App $app`: https://$name"
}
```

## Local Development (Codespace)

All demo apps can be built and run locally without Azure:

```powershell
# Example: Run TypeScript demo app
Set-Location cq-demo-app-001
docker build -t cq-demo-app-001 .
docker run -p 3000:3000 cq-demo-app-001
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| OIDC login fails | Verify environment name is `prod` (not `production`) |
| ACR name collision | Bicep uses `uniqueString()` — verify resource group name |
| Docker build fails | Ensure Docker Desktop is running |
| Coverage tools not found | Install with `pip install lizard` |
