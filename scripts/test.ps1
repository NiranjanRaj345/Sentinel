Write-Host "Running Go tests..."

Push-Location services/node-agent

go test ./...

$exitCode = $LASTEXITCODE

Pop-Location

if ($exitCode -ne 0) {
    exit $exitCode
}

Write-Host "All tests passed."