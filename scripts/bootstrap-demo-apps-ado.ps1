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
    [string]$AdoProject = "CodeQuality"
)

$ErrorActionPreference = "Stop"

$AdoOrgUrl = "https://dev.azure.com/$AdoOrg"
$ScannerRepo = "code-quality-scan-demo-app"
$DemoApps = @(
    @{ Name = "cq-demo-app-001"; Dir = "cq-demo-app-001"; Flavor = "javascript" },
    @{ Name = "cq-demo-app-002"; Dir = "cq-demo-app-002"; Flavor = "python" },
    @{ Name = "cq-demo-app-003"; Dir = "cq-demo-app-003"; Flavor = "dotnet" },
    @{ Name = "cq-demo-app-004"; Dir = "cq-demo-app-004"; Flavor = "java" },
    @{ Name = "cq-demo-app-005"; Dir = "cq-demo-app-005"; Flavor = "go" }
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
$RootDir = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path

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

    # Sync content — always push the latest template content to the app repo.
    # Compares a tree hash of the local template against the remote HEAD to
    # skip the push when content is already up to date.
    Write-Host "Syncing content to '$repoName'..." -ForegroundColor Green
    $repoUrl = "https://$AdoOrg@dev.azure.com/$AdoOrg/$AdoProject/_git/$repoName"

    # Export only tracked files (respects .gitignore) via git archive
    $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "ado-push-$repoName-$PID"
    $archiveFile = "$tempDir.zip"
    try {
        Push-Location $RootDir
        git archive --format=zip --output="$archiveFile" "HEAD:$appDir"
        Pop-Location

        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        Expand-Archive -Path $archiveFile -DestinationPath $tempDir -Force
        Remove-Item $archiveFile -Force

        Push-Location $tempDir
        git init | Out-Null
        git branch -M main
        git add -A

        # Check if anything differs from what's already in the remote
        $localTreeHash = git write-tree
        $needsPush = $true

        $defaultBranch = az repos show --repository $repoName --org $AdoOrgUrl --project $AdoProject --query defaultBranch -o tsv 2>$null
        if ($defaultBranch -and $defaultBranch -ne "None") {
            # Repo has content — fetch remote and compare tree hashes
            git remote add origin $repoUrl 2>$null
            git fetch origin main --depth=1 2>$null
            if ($LASTEXITCODE -eq 0) {
                $remoteTreeHash = git rev-parse "origin/main^{tree}" 2>$null
                if ($localTreeHash -eq $remoteTreeHash) {
                    $needsPush = $false
                    Write-Host "Repository '$repoName' content is up to date — skipping push." -ForegroundColor Yellow
                }
            }
        }

        if ($needsPush) {
            git commit -m "feat: sync scaffold for $repoName" --allow-empty | Out-Null
            if (-not (git remote | Select-String -SimpleMatch "origin")) {
                git remote add origin $repoUrl
            }
            git push -u origin main --force 2>&1
            if ($LASTEXITCODE -ne 0) {
                Write-Warning "Failed to push to '$repoName'. Verify git credentials for Azure DevOps."
            } else {
                Write-Host "Content pushed to '$repoName'." -ForegroundColor Green
            }
        }
        Pop-Location
    } finally {
        if (Test-Path $tempDir) { Remove-Item -Recurse -Force $tempDir }
        if (Test-Path $archiveFile) { Remove-Item -Force $archiveFile }
    }

    Write-Host "✅ ADO repo '$repoName' ready." -ForegroundColor Green
}

# ── Step 3b: Enable Advanced Security on all repos ──
Write-Host ""
Write-Host "--- Enabling Advanced Security on repositories ---" -ForegroundColor Cyan

$projectId = az devops project show --project $AdoProject --org $AdoOrgUrl --query id -o tsv

# Acquire an ADO-scoped bearer token for REST calls
$adoToken = az account get-access-token --resource "499b84ac-1321-427f-aa17-267ca6975798" --query accessToken -o tsv
$adoHeaders = @{ Authorization = "Bearer $adoToken"; "Content-Type" = "application/json" }
$AdvSecBaseUrl = "https://advsec.dev.azure.com/$AdoOrg/$projectId/_apis/management/repositories"

$allRepos = @($ScannerRepo) + ($DemoApps | ForEach-Object { $_.Name })
foreach ($repo in $allRepos) {
    $repoId = az repos show --repository $repo --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null
    if (-not $repoId) {
        Write-Host "  Repository '$repo' not found — skipping." -ForegroundColor Yellow
        continue
    }

    # Check current Advanced Security status via advsec.dev.azure.com
    $enablementUrl = "$AdvSecBaseUrl/$repoId/enablement?api-version=7.2-preview.1"
    try {
        $repoSettings = Invoke-RestMethod -Uri $enablementUrl -Method Get -Headers $adoHeaders
    } catch {
        $repoSettings = $null
    }
    if ($repoSettings -and $repoSettings.advSecEnabled -eq $true) {
        Write-Host "  Advanced Security already enabled on '$repo' — skipping." -ForegroundColor Yellow
    } else {
        Write-Host "  Enabling Advanced Security on '$repo'..." -ForegroundColor Green
        try {
            Invoke-RestMethod -Uri $enablementUrl -Method Patch -Headers $adoHeaders -Body '{"advSecEnabled": true}' -ContentType "application/json" | Out-Null
            Write-Host "  Advanced Security enabled on '$repo'." -ForegroundColor Green
        } catch {
            Write-Warning "  Failed to enable Advanced Security on '$repo': $($_.Exception.Message)"
        }
    }
}

# ── Step 4: Create pipeline definitions (scanner repo) ──
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

# ── Step 5: Create per-app scan pipelines (app repos) ──
Write-Host ""
Write-Host "--- Creating per-app scan pipelines ---" -ForegroundColor Cyan

foreach ($app in $DemoApps) {
    $repoName = $app.Name
    $flavor = $app.Flavor
    $pipelineName = "Code Quality Scan - $repoName"
    $existingPipeline = az pipelines show --name $pipelineName --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null

    if ($existingPipeline) {
        Write-Host "Pipeline '$pipelineName' already exists — skipping creation." -ForegroundColor Yellow
    } else {
        Write-Host "Creating pipeline '$pipelineName' in repo '$repoName'..." -ForegroundColor Green
        az pipelines create `
            --name $pipelineName `
            --repository $repoName `
            --repository-type tfsgit `
            --branch main `
            --yml-path ".azuredevops/pipelines/code-quality-scan.yml" `
            --org $AdoOrgUrl `
            --project $AdoProject `
            --skip-first-run
        Write-Host "Pipeline '$pipelineName' created." -ForegroundColor Green
    }

    # Set the MegaLinter flavor variable on the pipeline (idempotent — overwrites)
    $pipelineId = az pipelines show --name $pipelineName --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null
    if ($pipelineId) {
        Write-Host "  Setting megalinterFlavor='$flavor' on '$pipelineName'..." -ForegroundColor Gray
        $defUrl = "$AdoOrgUrl/$projectId/_apis/build/definitions/${pipelineId}?api-version=7.1"
        try {
            $def = Invoke-RestMethod -Uri $defUrl -Method Get -Headers $adoHeaders
            # Set variables property via Add-Member to handle missing property on PSObject
            $varObj = @{ megalinterFlavor = @{ value = $flavor; allowOverride = $true } }
            $def | Add-Member -NotePropertyName variables -NotePropertyValue $varObj -Force
            $updatedBody = $def | ConvertTo-Json -Depth 50 -Compress
            $null = Invoke-RestMethod -Uri $defUrl -Method Put -Headers $adoHeaders -Body $updatedBody -ContentType "application/json"
            Write-Host "  megalinterFlavor='$flavor' set." -ForegroundColor Green
        } catch {
            Write-Warning "  Failed to set megalinterFlavor on '$pipelineName': $($_.Exception.Message)"
        }
    }
}

# ── Step 6: Authorize variable groups and grant pipeline permissions ──
Write-Host ""
Write-Host "--- Authorizing variable groups and pipeline permissions ---" -ForegroundColor Cyan

# Acquire an ADO-scoped bearer token for REST calls
$adoToken = az account get-access-token --resource "499b84ac-1321-427f-aa17-267ca6975798" --query accessToken -o tsv
$adoHeaders = @{ Authorization = "Bearer $adoToken"; "Content-Type" = "application/json" }
$apiVersion = "api-version=7.1-preview.1"

function Invoke-AdoPatch {
    param([string]$Url, [hashtable]$Body)
    $json = $Body | ConvertTo-Json -Depth 10 -Compress
    try {
        $resp = Invoke-RestMethod -Uri $Url -Method Patch -Headers $adoHeaders -Body $json -ContentType "application/json"
        return $resp
    } catch {
        Write-Warning "  REST PATCH failed: $($_.Exception.Message)"
        return $null
    }
}

function Invoke-AdoGet {
    param([string]$Url)
    try {
        return Invoke-RestMethod -Uri $Url -Method Get -Headers $adoHeaders
    } catch {
        return $null
    }
}

# Authorize variable groups for all pipelines
foreach ($vg in @($vgName, $oidcVgName)) {
    $vgId = az pipelines variable-group list --org $AdoOrgUrl --project $AdoProject --query "[?name=='$vg'].id" -o tsv 2>$null
    if (-not $vgId) { continue }

    $url = "$AdoOrgUrl/$projectId/_apis/pipelines/pipelinepermissions/variablegroup/$vgId`?$apiVersion"
    $existing = Invoke-AdoGet -Url $url
    if ($existing -and $existing.allPipelines -and $existing.allPipelines.authorized -eq $true) {
        Write-Host "Variable group '$vg' already authorized for all pipelines — skipping." -ForegroundColor Yellow
    } else {
        Write-Host "Authorizing variable group '$vg' for all pipelines..." -ForegroundColor Green
        $result = Invoke-AdoPatch -Url $url -Body @{
            resource     = @{ id = $vgId; type = "variablegroup" }
            allPipelines = @{ authorized = $true }
            pipelines    = @()
        }
        if ($result) {
            Write-Host "Variable group '$vg' authorized." -ForegroundColor Green
        } else {
            Write-Warning "Failed to authorize variable group '$vg'. Authorize manually in Pipelines > Library."
        }
    }
}

# Authorize the scanner repo for all pipelines (per-app pipelines check it out)
$scannerRepoId = az repos show --repository $ScannerRepo --org $AdoOrgUrl --project $AdoProject --query id -o tsv 2>$null
if ($scannerRepoId) {
    $url = "$AdoOrgUrl/$projectId/_apis/pipelines/pipelinepermissions/repository/$projectId.$scannerRepoId`?$apiVersion"
    $existing = Invoke-AdoGet -Url $url
    if ($existing -and $existing.allPipelines -and $existing.allPipelines.authorized -eq $true) {
        Write-Host "Scanner repo already authorized for all pipelines — skipping." -ForegroundColor Yellow
    } else {
        Write-Host "Authorizing scanner repo for all pipelines..." -ForegroundColor Green
        $result = Invoke-AdoPatch -Url $url -Body @{
            resource     = @{ id = "$projectId.$scannerRepoId"; type = "repository" }
            allPipelines = @{ authorized = $true }
            pipelines    = @()
        }
        if ($result) {
            Write-Host "Scanner repo authorized for all pipelines." -ForegroundColor Green
        } else {
            Write-Warning "Failed to authorize scanner repo. Authorize manually."
        }
    }
}

Write-Host ""
Write-Host "=== ADO Bootstrap Complete ===" -ForegroundColor Cyan
Write-Host "All ADO repos, variable groups, pipelines, and permissions have been configured."
