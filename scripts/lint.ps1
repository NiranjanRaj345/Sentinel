Write-Host "Running go vet..."

Push-Location services/node-agent

go vet ./...

$exitCode = $LASTEXITCODE

Pop-Location

if ($exitCode -ne 0) {
    exit $exitCode
}

Write-Host "Lint completed successfully."