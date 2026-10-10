# Launch the Try Arch setup with an explicitly verified vanilla guest release.
param(
    [Parameter(Mandatory=$true)][string]$Release,
    [Parameter(Mandatory=$true)][string]$SumsSHA256,
    [string]$Dir = "$env:LOCALAPPDATA\TryArch"
)
$ErrorActionPreference = 'Stop'
if ($SumsSHA256 -notmatch '^[0-9a-fA-F]{64}$') { throw 'A trusted SHA256SUMS digest is required.' }
$launcher = Join-Path $PSScriptRoot '../TryArch.exe'
if (-not (Test-Path -LiteralPath $launcher)) { throw 'Build TryArch.exe first; see README.md.' }
& $launcher -dir $Dir -release $Release -sums-sha256 $SumsSHA256
