<#
.SYNOPSIS
    Deploy Code Quality Report PBIP to a Power BI workspace using FabricPS-PBIP module.

.DESCRIPTION
    This script deploys the Code Quality Report (PBIP format) to a Power BI workspace.
    It uses the FabricPS-PBIP module for Fabric/Power BI Git integration.

    Prerequisites:
    - PowerShell 7+
    - FabricPS-PBIP module (Install-Module FabricPS-PBIP)
    - Power BI Pro or Premium Per User license
    - Azure AD authentication configured

.PARAMETER WorkspaceName
    Name of the target Power BI workspace. Created if it does not exist.

.PARAMETER WorkspaceId
    Optional. GUID of an existing workspace. Takes precedence over WorkspaceName.

.PARAMETER PbipPath
    Path to the PBIP file. Defaults to the CodeQualityReport.pbip in the same directory.

.PARAMETER Environment
    Deployment environment (dev, staging, prod). Affects data source parameters.

.EXAMPLE
    .\deploy.ps1 -WorkspaceName "Code Quality Reports - Dev"

.EXAMPLE
    .\deploy.ps1 -WorkspaceName "Code Quality Reports - Prod" -Environment prod
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceName,

    [Parameter(Mandatory = $false)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $false)]
    [string]$PbipPath = (Join-Path $PSScriptRoot ".." "CodeQualityReport.pbip"),

    [Parameter(Mandatory = $false)]
    [ValidateSet("dev", "staging", "prod")]
    [string]$Environment = "dev"
)

$ErrorActionPreference = "Stop"

# ─── Module Check ───────────────────────────────────────────────────────────────

Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Code Quality Report — Power BI Deployment" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

if (-not (Get-Module -ListAvailable -Name "FabricPS-PBIP")) {
    Write-Host "Installing FabricPS-PBIP module..." -ForegroundColor Yellow
    Install-Module -Name "FabricPS-PBIP" -Scope CurrentUser -Force -AllowClobber
}
Import-Module FabricPS-PBIP -Force

# ─── Authentication ─────────────────────────────────────────────────────────────

Write-Host "Authenticating to Power BI service..." -ForegroundColor Yellow
$token = Get-FabricAuthToken

if (-not $token) {
    Write-Error "Failed to authenticate to Power BI. Ensure you have a valid Azure AD session."
    exit 1
}
Write-Host "✅ Authenticated successfully." -ForegroundColor Green

# ─── Workspace Resolution ──────────────────────────────────────────────────────

if ($WorkspaceId) {
    Write-Host "Using provided workspace ID: $WorkspaceId" -ForegroundColor Yellow
} else {
    Write-Host "Resolving workspace '$WorkspaceName'..." -ForegroundColor Yellow

    $workspace = Get-FabricWorkspace -Name $WorkspaceName -ErrorAction SilentlyContinue

    if (-not $workspace) {
        Write-Host "Workspace '$WorkspaceName' not found. Creating..." -ForegroundColor Yellow
        $workspace = New-FabricWorkspace -Name $WorkspaceName
        Write-Host "✅ Workspace created: $($workspace.id)" -ForegroundColor Green
    } else {
        Write-Host "✅ Workspace found: $($workspace.id)" -ForegroundColor Green
    }

    $WorkspaceId = $workspace.id
}

# ─── Validate PBIP Path ────────────────────────────────────────────────────────

$resolvedPbipPath = Resolve-Path $PbipPath -ErrorAction SilentlyContinue

if (-not $resolvedPbipPath) {
    Write-Error "PBIP file not found at: $PbipPath"
    exit 1
}

Write-Host "PBIP path: $resolvedPbipPath" -ForegroundColor Yellow

# ─── Deploy Semantic Model ──────────────────────────────────────────────────────

$semanticModelPath = Join-Path (Split-Path $resolvedPbipPath) "CodeQualityReport.SemanticModel"

if (Test-Path $semanticModelPath) {
    Write-Host ""
    Write-Host "Deploying semantic model..." -ForegroundColor Yellow

    $semanticModel = Import-FabricItem `
        -WorkspaceId $WorkspaceId `
        -Path $semanticModelPath `
        -Type "SemanticModel"

    Write-Host "✅ Semantic model deployed: $($semanticModel.displayName)" -ForegroundColor Green
} else {
    Write-Error "Semantic model directory not found at: $semanticModelPath"
    exit 1
}

# ─── Deploy Report ──────────────────────────────────────────────────────────────

$reportPath = Join-Path (Split-Path $resolvedPbipPath) "CodeQualityReport.Report"

if (Test-Path $reportPath) {
    Write-Host ""
    Write-Host "Deploying report..." -ForegroundColor Yellow

    $report = Import-FabricItem `
        -WorkspaceId $WorkspaceId `
        -Path $reportPath `
        -Type "Report"

    Write-Host "✅ Report deployed: $($report.displayName)" -ForegroundColor Green
} else {
    Write-Error "Report directory not found at: $reportPath"
    exit 1
}

# ─── Configure Data Source Parameters ───────────────────────────────────────────

Write-Host ""
Write-Host "Configuring data source parameters for '$Environment' environment..." -ForegroundColor Yellow

$setupParamsScript = Join-Path $PSScriptRoot "setup-parameters.ps1"
if (Test-Path $setupParamsScript) {
    & $setupParamsScript -WorkspaceId $WorkspaceId -SemanticModelId $semanticModel.id -Environment $Environment
} else {
    Write-Host "⚠️  setup-parameters.ps1 not found — skipping parameter configuration." -ForegroundColor Yellow
    Write-Host "   Run setup-parameters.ps1 manually to configure data source." -ForegroundColor Yellow
}

# ─── Refresh Dataset ───────────────────────────────────────────────────────────

Write-Host ""
Write-Host "Triggering dataset refresh..." -ForegroundColor Yellow

try {
    Invoke-FabricItemRefresh -WorkspaceId $WorkspaceId -ItemId $semanticModel.id
    Write-Host "✅ Dataset refresh initiated." -ForegroundColor Green
} catch {
    Write-Host "⚠️  Dataset refresh failed (data source may need credentials configured)." -ForegroundColor Yellow
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host "   Configure data source credentials in Power BI service, then refresh manually." -ForegroundColor Yellow
}

# ─── Summary ────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Green
Write-Host "  Deployment Complete" -ForegroundColor Green
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Green
Write-Host ""
Write-Host "  Workspace:      $WorkspaceName" -ForegroundColor White
Write-Host "  Workspace ID:   $WorkspaceId" -ForegroundColor White
Write-Host "  Semantic Model: $($semanticModel.displayName)" -ForegroundColor White
Write-Host "  Report:         $($report.displayName)" -ForegroundColor White
Write-Host "  Environment:    $Environment" -ForegroundColor White
Write-Host ""
Write-Host "  Next Steps:" -ForegroundColor Yellow
Write-Host "    1. Open Power BI service and navigate to the workspace." -ForegroundColor White
Write-Host "    2. Configure OAuth credentials for ADLS Gen2 data source." -ForegroundColor White
Write-Host "    3. Refresh the dataset to load scan results." -ForegroundColor White
Write-Host ""
