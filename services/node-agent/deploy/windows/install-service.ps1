param(
    [string]$BinaryPath = "$(Split-Path -Parent $MyInvocation.MyCommand.Path)\sentinel-service.exe",
    [string]$ServiceName = "sentinel-agent",
    [string]$DisplayName = "Sentinel Node Agent",
    [string]$Description = "Sentinel node monitoring and recovery agent"
)

if (-not (Test-Path $BinaryPath)) {
    Write-Error "Binary not found at $BinaryPath"
    exit 1
}

$BinaryPath = Resolve-Path $BinaryPath

if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Stopping existing service..."
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $ServiceName | Out-Null
}

Write-Host "Installing service..."
sc.exe create $ServiceName binPath= "`"$BinaryPath`"" start= auto DisplayName= $DisplayName | Out-Null
sc.exe description $ServiceName $Description | Out-Null

Write-Host "Starting service..."
Start-Service -Name $ServiceName

Write-Host "Service installed and started."
Write-Host "View logs: Get-Content -Path 'logs\sentinel.log' -Wait"
