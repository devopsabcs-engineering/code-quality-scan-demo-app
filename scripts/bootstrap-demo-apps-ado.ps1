<#
.SYNOPSIS
    Bootstraps Azure DevOps repos, variable groups, and pipelines for Code Quality.

.DESCRIPTION
    Creates Azure DevOps repositories, pushes content, creates variable groups,
    service connections, and pipeline definitions. All operations are idempotent.

.PARAMETER AdoOrg
    Azure DevOps organization name.

.PARAMETER AdoProject
    Azure DevOps project name.

.PARAMETER ClientId
    Azure AD client ID from OIDC setup.

.PARAMETER TenantId
    Azure AD tenant ID.

.PARAMETER SubscriptionId
    Azure subscription ID.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ClientId,

    [Parameter(Mandatory = $true)]
    [string]$TenantId,

    [Parameter(Mandatory = $true)]
    [string]$SubscriptionId,

    [Parameter(Mandatory = $false)]
    [string]$AdoOrg = "MngEnvMCAP675646",

    [Parameter(Mandatory = $false)]
    [string]$AdoProject = "Agentic Accelerator Framework"
)

$ErrorActionPreference = "Stop"

$AdoOrgUrl = "https://dev.azure.com/$AdoOrg"
$ScannerRepo = "code-quality-scan-demo-app"
$DemoApps = @(
    @{ Name = "cq-demo-app-001"; Dir = "cq-demo-app-001" },
    @{ Name = "cq-demo-app-002"; Dir = "cq-demo-app-002" },
    @{ Name = "cq-demo-app-003"; Dir = "cq-demo-app-003" },
    @{ Name = "cq-demo-app-004"; Dir = "cq-demo-app-004" },
    @{ Name = "cq-demo-app-005"; Dir = "cq-demo-app-005" }
)

Write-Host "=== Code Quality ADO Bootstrap ===" -ForegroundColor Cyan
Write-Host "ADO Org: $AdoOrg"
Write-Host "ADO Project: $AdoProject"
Write-Host ""

# ── Step 1: Create variable group ──
$vgName = "code-quality-common"
$existingVg = az pipelines variable-group list --org $AdoOrgUrl --project $AdoProject --query "[?name=='$vgName'].id" -o tsv 2>$null

if ($existingVg) {
    Write-Host "Variable group '$vgName' already exists (id: $existingVg) — skipping." -ForegroundColor Yellow
} else {
    Write-Host "Creating variable group '$vgName'..." -ForegroundColor Green
    az pipelines variable-group create `
        --org $AdoOrgUrl `
        --project $AdoProject `
        --name $vgName `
        --variables `
            "location=canadacentral" `
            "serviceConnection=code-quality-scan-demo-app-ado-sc"
    Write-Host "Variable group '$vgName' created." -ForegroundColor Green
}

# ── Step 2: Create OIDC variable group ──
$oidcVgName = "code-quality-oidc"
$existingOidcVg = az pipelines variable-group list --org $AdoOrgUrl --project $AdoProject --query "[?name=='$oidcVgName'].id" -o tsv 2>$null

if ($existingOidcVg) {
    Write-Host "Variable group '$oidcVgName' already exists — skipping." -ForegroundColor Yellow
} else {
    Write-Host "Creating variable group '$oidcVgName'..." -ForegroundColor Green
    az pipelines variable-group create `
        --org $AdoOrgUrl `
        --project $AdoProject `
        --name $oidcVgName `
        --variables `
            "clientId=$ClientId" `
            "tenantId=$TenantId" `
            "subscriptionId=$SubscriptionId"
    Write-Host "Variable group '$oidcVgName' created." -ForegroundColor Green
}

# ── Step 3: Create ADO repos and push content ──
foreach ($app in $DemoApps) {
    $repoName = $app.Name
    $appDir = $app.Dir

    Write-Host ""
    Write-Host "--- Processing ADO repo: $repoName ---" -ForegroundColor Cyan

    # Guard: Check if repo exists
    $existingRepo = az repos show --repository $repoName --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null
    if ($existingRepo) {
        Write-Host "ADO repository '$repoName' already exists — skipping creation." -ForegroundColor Yellow
    } else {
        Write-Host "Creating ADO repository '$repoName'..." -ForegroundColor Green
        az repos create --name $repoName --org $AdoOrgUrl --project $AdoProject
    }

    Write-Host "✅ ADO repo '$repoName' ready." -ForegroundColor Green
}

# ── Step 4: Create pipeline definitions ──
$pipelines = @(
    @{ Name = "Code Quality Scan"; YamlPath = ".azuredevops/pipelines/code-quality-scan.yml" },
    @{ Name = "Code Quality Lint Gate"; YamlPath = ".azuredevops/pipelines/code-quality-lint-gate.yml" },
    @{ Name = "Deploy All"; YamlPath = ".azuredevops/pipelines/deploy-all.yml" },
    @{ Name = "Teardown All"; YamlPath = ".azuredevops/pipelines/teardown-all.yml" },
    @{ Name = "Scan and Store"; YamlPath = ".azuredevops/pipelines/scan-and-store.yml" }
)

foreach ($pipeline in $pipelines) {
    $pipelineName = $pipeline.Name
    $existingPipeline = az pipelines show --name $pipelineName --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null

    if ($existingPipeline) {
        Write-Host "Pipeline '$pipelineName' already exists — skipping." -ForegroundColor Yellow
    } else {
        Write-Host "Creating pipeline '$pipelineName'..." -ForegroundColor Green
        az pipelines create `
            --name $pipelineName `
            --repository $ScannerRepo `
            --repository-type tfsgit `
            --branch main `
            --yml-path $pipeline.YamlPath `
            --org $AdoOrgUrl `
            --project $AdoProject `
            --skip-first-run
        Write-Host "Pipeline '$pipelineName' created." -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "=== ADO Bootstrap Complete ===" -ForegroundColor Cyan
Write-Host "All ADO repos, variable groups, and pipelines have been configured."
