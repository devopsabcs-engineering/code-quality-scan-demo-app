<#
.SYNOPSIS
    Configure data source parameters for the Code Quality Report semantic model.

.DESCRIPTION
    Updates the ADLS Gen2 connection parameters (storage account name, container name)
    for the Code Quality Report semantic model in a Power BI workspace.

    This script is called automatically by deploy.ps1, but can also be run
    independently to reconfigure data source parameters.

    Prerequisites:
    - PowerShell 7+
    - FabricPS-PBIP module
    - Power BI Pro or Premium Per User license

.PARAMETER WorkspaceId
    GUID of the Power BI workspace containing the semantic model.

.PARAMETER SemanticModelId
    GUID of the semantic model to configure. If not provided, searches by name.

.PARAMETER Environment
    Deployment environment. Controls which storage account is used:
    - dev:     stcqscandev
    - staging: stcqscanstaging
    - prod:    stcqscan

.PARAMETER StorageAccountName
    Optional override for the ADLS Gen2 storage account name.
    Takes precedence over the environment-based default.

.PARAMETER ContainerName
    Optional override for the ADLS Gen2 container name.
    Default: code-quality-results

.EXAMPLE
    .\setup-parameters.ps1 -WorkspaceId "abc-123" -Environment prod

.EXAMPLE
    .\setup-parameters.ps1 -WorkspaceId "abc-123" -StorageAccountName "mystorageacct" -ContainerName "scan-data"
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $false)]
    [string]$SemanticModelId,

    [Parameter(Mandatory = $false)]
    [ValidateSet("dev", "staging", "prod")]
    [string]$Environment = "dev",

    [Parameter(Mandatory = $false)]
    [string]$StorageAccountName,

    [Parameter(Mandatory = $false)]
    [string]$ContainerName = "code-quality-results"
)

$ErrorActionPreference = "Stop"

# ─── Module Check ───────────────────────────────────────────────────────────────

if (-not (Get-Module -ListAvailable -Name "FabricPS-PBIP")) {
    Write-Error "FabricPS-PBIP module not installed. Run: Install-Module FabricPS-PBIP -Scope CurrentUser"
    exit 1
}
Import-Module FabricPS-PBIP -Force

# ─── Resolve Storage Account ───────────────────────────────────────────────────

if (-not $StorageAccountName) {
    $StorageAccountName = switch ($Environment) {
        "dev"     { "stcqscandev" }
        "staging" { "stcqscanstaging" }
        "prod"    { "stcqscan" }
    }
}

Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Code Quality Report — Data Source Configuration" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Environment:     $Environment" -ForegroundColor White
Write-Host "  Storage Account: $StorageAccountName" -ForegroundColor White
Write-Host "  DFS Endpoint:    https://$StorageAccountName.dfs.core.windows.net" -ForegroundColor White
Write-Host "  Container:       $ContainerName" -ForegroundColor White
Write-Host ""

# ─── Resolve Semantic Model ────────────────────────────────────────────────────

if (-not $SemanticModelId) {
    Write-Host "Searching for 'CodeQualityReport' semantic model..." -ForegroundColor Yellow

    $items = Get-FabricWorkspaceItem -WorkspaceId $WorkspaceId -Type "SemanticModel"
    $model = $items | Where-Object { $_.displayName -eq "CodeQualityReport" } | Select-Object -First 1

    if (-not $model) {
        Write-Error "Semantic model 'CodeQualityReport' not found in workspace $WorkspaceId. Deploy first using deploy.ps1."
        exit 1
    }

    $SemanticModelId = $model.id
    Write-Host "✅ Found semantic model: $SemanticModelId" -ForegroundColor Green
} else {
    Write-Host "Using provided semantic model ID: $SemanticModelId" -ForegroundColor Yellow
}

# ─── Update Parameters ─────────────────────────────────────────────────────────

Write-Host ""
Write-Host "Updating data source parameters..." -ForegroundColor Yellow

$parameters = @{
    updateDetails = @(
        @{
            name  = "StorageAccountName"
            newValue = $StorageAccountName
        },
        @{
            name  = "ContainerName"
            newValue = $ContainerName
        }
    )
}

try {
    # Use Power BI REST API to update parameters
    $apiUrl = "https://api.powerbi.com/v1.0/myorg/groups/$WorkspaceId/datasets/$SemanticModelId/Default.UpdateParameters"

    $token = Get-FabricAuthToken
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type"  = "application/json"
    }

    $body = $parameters | ConvertTo-Json -Depth 10

    Invoke-RestMethod -Uri $apiUrl -Method Post -Headers $headers -Body $body

    Write-Host "✅ Parameters updated successfully." -ForegroundColor Green
} catch {
    Write-Host "⚠️  Parameter update via REST API failed." -ForegroundColor Yellow
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "   Manual Configuration:" -ForegroundColor Yellow
    Write-Host "   1. Open Power BI service → Workspace → Settings → Parameters" -ForegroundColor White
    Write-Host "   2. Set StorageAccountName = $StorageAccountName" -ForegroundColor White
    Write-Host "   3. Set ContainerName = $ContainerName" -ForegroundColor White
    Write-Host ""
}

# ─── Configure Credentials ─────────────────────────────────────────────────────

Write-Host ""
Write-Host "Credential Configuration:" -ForegroundColor Yellow
Write-Host "  The ADLS Gen2 data source uses OAuth (Organizational Account)." -ForegroundColor White
Write-Host "  After deployment, configure credentials in Power BI service:" -ForegroundColor White
Write-Host ""
Write-Host "  1. Navigate to Workspace → Settings → Data source credentials" -ForegroundColor White
Write-Host "  2. Find the ADLS Gen2 data source" -ForegroundColor White
Write-Host "  3. Click 'Edit credentials'" -ForegroundColor White
Write-Host "  4. Authentication method: OAuth2" -ForegroundColor White
Write-Host "  5. Privacy level: Organizational" -ForegroundColor White
Write-Host "  6. Sign in with your Azure AD account" -ForegroundColor White
Write-Host ""

# ─── Validate Connection ───────────────────────────────────────────────────────

Write-Host "Data Source Connection String:" -ForegroundColor Cyan
Write-Host "  AzureStorage.DataLake(`"https://$StorageAccountName.dfs.core.windows.net/$ContainerName`")" -ForegroundColor White
Write-Host ""
Write-Host "Expected ADLS Gen2 Path Structure:" -ForegroundColor Cyan
Write-Host "  $ContainerName/" -ForegroundColor White
Write-Host "    ├── {yyyy}/{MM}/{dd}/{appId}-{tool}.json    (scan results)" -ForegroundColor Gray
Write-Host "    └── metadata/" -ForegroundColor White
Write-Host "        ├── repositories.json                    (repo dimension)" -ForegroundColor Gray
Write-Host "        ├── scan-tools.json                      (tool dimension)" -ForegroundColor Gray
Write-Host "        └── languages.json                       (language dimension)" -ForegroundColor Gray
Write-Host ""

Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Green
Write-Host "  Configuration Complete" -ForegroundColor Green
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Green
