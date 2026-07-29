Write-Host "Formatting Go source..."

Push-Location services/node-agent

go fmt ./...

$exitCode = $LASTEXITCODE

Pop-Location

if ($exitCode -ne 0) {
    exit $exitCode
}

Write-Host "Formatting completed successfully."