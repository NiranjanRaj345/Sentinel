Write-Host "Building Sentinel Node Agent..."

Push-Location services/node-agent

go build ./...

$exitCode = $LASTEXITCODE

Pop-Location

if ($exitCode -ne 0) {
    exit $exitCode
}

Write-Host "Build successful."