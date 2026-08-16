# Build eudic-mcp-go for Windows + Linux (Cursor Remote SSH needs linux/*).
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root
New-Item -ItemType Directory -Force -Path dist/linux-amd64, dist/linux-arm64, dist/windows-amd64 | Out-Null
$env:CGO_ENABLED = "0"

Write-Host "-> linux/amd64"
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "dist/linux-amd64/eudic-mcp-go" .

Write-Host "-> linux/arm64"
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go build -ldflags="-s -w" -o "dist/linux-arm64/eudic-mcp-go" .

Write-Host "-> windows/amd64"
$env:GOOS = "windows"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "dist/windows-amd64/eudic-mcp-go.exe" .

Write-Host "Done. See dist/"
Get-ChildItem dist -Recurse -File | Format-Table FullName, Length
