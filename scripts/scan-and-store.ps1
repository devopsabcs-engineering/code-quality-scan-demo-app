<#
.SYNOPSIS
    Parses SARIF scan results and uploads to ADLS Gen2.

.DESCRIPTION
    Reads SARIF files from a directory, extracts findings,
    and uploads to Azure Data Lake Storage Gen2 with date-partitioned paths.
    Path pattern: {yyyy}/{MM}/{dd}/{appId}-{tool}.json

.PARAMETER StorageAccountName
    ADLS Gen2 storage account name.

.PARAMETER ContainerName
    Blob container name.

.PARAMETER SarifDir
    Directory containing SARIF files.

.PARAMETER AppId
    Demo app identifier (e.g., '001').
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$StorageAccountName,

    [Parameter(Mandatory = $true)]
    [string]$ContainerName,

    [Parameter(Mandatory = $true)]
    [string]$SarifDir,

    [Parameter(Mandatory = $true)]
    [string]$AppId
)

$ErrorActionPreference = "Stop"

$datePath = Get-Date -Format "yyyy/MM/dd"
$timestamp = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

Write-Host "=== Scan and Store ===" -ForegroundColor Cyan
Write-Host "Storage Account: $StorageAccountName"
Write-Host "Container: $ContainerName"
Write-Host "SARIF Directory: $SarifDir"
Write-Host "App ID: $AppId"
Write-Host "Date Path: $datePath"
Write-Host ""

# ── Find all SARIF files ──
$sarifFiles = Get-ChildItem -Path $SarifDir -Filter "*.sarif" -Recurse

if ($sarifFiles.Count -eq 0) {
    Write-Host "No SARIF files found in $SarifDir" -ForegroundColor Yellow
    exit 0
}

Write-Host "Found $($sarifFiles.Count) SARIF file(s)." -ForegroundColor Green

foreach ($sarifFile in $sarifFiles) {
    Write-Host ""
    Write-Host "Processing: $($sarifFile.Name)" -ForegroundColor Cyan

    # Parse SARIF
    $sarif = Get-Content -Path $sarifFile.FullName -Raw | ConvertFrom-Json

    foreach ($run in $sarif.runs) {
        $toolName = $run.tool.driver.name
        $findings = $run.results

        Write-Host "  Tool: $toolName — $($findings.Count) findings"

        # Build flattened records for Power BI
        $records = @()
        foreach ($finding in $findings) {
            $location = $finding.locations[0].physicalLocation
            $record = @{
                finding_id   = [guid]::NewGuid().ToString()
                app_id       = "cq-demo-app-$AppId"
                rule_id      = $finding.ruleId
                level        = $finding.level
                message      = $finding.message.text
                file_path    = $location.artifactLocation.uri
                start_line   = $location.region.startLine
                tool_name    = $toolName
                scan_date    = $timestamp
                category     = if ($finding.properties.tags) { $finding.properties.tags -join "," } else { "code-quality" }
            }
            $records += $record
        }

        # Upload to ADLS Gen2
        $blobName = "$datePath/$AppId-$($toolName.ToLower()).json"
        $jsonContent = $records | ConvertTo-Json -Depth 10

        $tempUploadFile = [System.IO.Path]::GetTempFileName()
        $jsonContent | Set-Content -Path $tempUploadFile -Encoding UTF8

        Write-Host "  Uploading to: $ContainerName/$blobName"
        az storage fs file upload `
            --account-name $StorageAccountName `
            --file-system $ContainerName `
            --path $blobName `
            --source $tempUploadFile `
            --overwrite true `
            --auth-mode login

        Remove-Item -Path $tempUploadFile -Force
        Write-Host "  ✅ Uploaded $($records.Count) records for $toolName" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "=== Upload Complete ===" -ForegroundColor Cyan
