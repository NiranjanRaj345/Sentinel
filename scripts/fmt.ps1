Write-Host "Formatting Go source..."

go fmt ./services/node-agent/...

if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

Write-Host "Formatting completed successfully."