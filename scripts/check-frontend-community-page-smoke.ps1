param(
    [int]$BackendPort = 19997,
    [int]$FrontendPort = 3001,
    [int]$TimeoutSeconds = 90
)

$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path -LiteralPath ".").Path
$scriptPath = Join-Path $repoRoot "scripts/frontend-community-page-smoke.cjs"

& node $scriptPath `
    --backend-port $BackendPort `
    --frontend-port $FrontendPort `
    --timeout-seconds $TimeoutSeconds

if ($LASTEXITCODE -ne 0) {
    throw "frontend community page smoke failed with exit code $LASTEXITCODE"
}
