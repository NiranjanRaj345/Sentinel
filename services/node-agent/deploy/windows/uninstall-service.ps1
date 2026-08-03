param(
    [string]$ServiceName = "sentinel-agent"
)

if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Stopping service..."
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    Write-Host "Removing service..."
    sc.exe delete $ServiceName | Out-Null
    Write-Host "Service removed."
} else {
    Write-Host "Service not found."
}
