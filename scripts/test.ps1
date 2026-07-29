Write-Host "Running Go tests..."

go test ./services/node-agent/...

if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

Write-Host "All tests passed."