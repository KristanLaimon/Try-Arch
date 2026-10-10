param(
    [Parameter(Mandatory)][string]$Artifacts,
    [Parameter(Mandatory)]
    [ValidatePattern('^https://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+/releases/download/v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?$')]
    [string]$ReleaseUrl
)
$ErrorActionPreference = 'Stop'
$repo = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$assets = (Resolve-Path -LiteralPath $Artifacts).Path
$manifest = Join-Path $assets 'SHA256SUMS'
$pin = (Get-FileHash -LiteralPath $manifest -Algorithm SHA256).Hash.ToLowerInvariant()
$versionSource = Get-Content -LiteralPath (Join-Path $repo 'app/update.go') -Raw
$version = [regex]::Match($versionSource, 'currentVersion\s*=\s*"([^"]+)"').Groups[1].Value
if (($ReleaseUrl -split '/')[-1] -cne $version) { throw "Release tag must match $version in app/update.go" }
# Embed the exact manifest; hashing the executable separately avoids a circular pin.
Copy-Item -LiteralPath $manifest -Destination (Join-Path $repo 'app/testdata/SHA256SUMS.arch') -Force
$pinSource = Join-Path $repo 'app/manifest.go'
$pinText = [IO.File]::ReadAllText($pinSource)
$pinPattern = '(defaultSumsSHA256\s*=\s*)"[0-9a-f]{64}"'
if ([regex]::Matches($pinText, $pinPattern).Count -ne 1) { throw 'Guest pin declaration is missing or ambiguous' }
$pinText = [regex]::Replace($pinText, $pinPattern, '${1}"' + $pin + '"')
[IO.File]::WriteAllText($pinSource, $pinText, [Text.UTF8Encoding]::new($false))
$pins = "-X main.defaultReleaseURL=$ReleaseUrl -X main.defaultSumsSHA256=$pin -X main.defaultRuntimeReleaseURL=$ReleaseUrl -X main.defaultRuntimeSumsSHA256=$pin"
Push-Location (Join-Path $repo 'app')
try {
    go test -ldflags $pins ./...
    if ($LASTEXITCODE -ne 0) { throw 'Launcher release tests failed' }
    $executable = Join-Path $assets 'TryArch.exe'
    go build -trimpath -ldflags "-H windowsgui -s -w $pins" -o $executable .
    if ($LASTEXITCODE -ne 0) { throw 'Launcher release build failed' }
    $info = (Get-Item -LiteralPath $executable).VersionInfo
    if ($info.ProductName -cne 'Try Arch' -or $info.ProductVersion -cne $version) { throw 'Executable metadata is stale' }
    $digest = (Get-FileHash -LiteralPath $executable -Algorithm SHA256).Hash.ToLowerInvariant()
    [IO.File]::WriteAllText("$executable.sha256", "$digest  TryArch.exe`n", [Text.UTF8Encoding]::new($false))
    Write-Host "Built $executable ($version) with guest/runtime manifest pin $pin"
} finally {
    Pop-Location
}
