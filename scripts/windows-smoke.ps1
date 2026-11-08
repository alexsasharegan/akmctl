<#
.SYNOPSIS
  Runs the real akmctl binary the way a user does, on a temporary repository.

.DESCRIPTION
  Builds nothing: pass the path of an akmctl binary (akmctl.exe on Windows).
  For each mode (git core.autocrlf true, then false) it creates a repository in a
  directory whose path has spaces, and drives init, new, set, section, claim,
  close (a file move that rewrites links) and check. It also writes while a
  reader holds a record open. Any wrong exit code or output fails the script
  (exit 1). It runs on Windows (CI, windows-latest) and, with pwsh installed,
  on macOS and Linux.

  Work item 261108-hwcq. What it covers: docs/architecture.md, "Platform coverage".

.EXAMPLE
  go build -o $env:TEMP\akmctl.exe .\cmd\akmctl
  pwsh -File scripts\windows-smoke.ps1 -Exe $env:TEMP\akmctl.exe
#>
param(
  [Parameter(Mandatory = $true)][string]$Exe,
  [string]$Work = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$Exe = (Resolve-Path -LiteralPath $Exe).Path
if (-not $Work) { $Work = Join-Path ([System.IO.Path]::GetTempPath()) ("akm smoke " + [guid]::NewGuid().ToString('N').Substring(0, 8)) }
New-Item -ItemType Directory -Force -Path $Work | Out-Null
$Work = (Resolve-Path -LiteralPath $Work).Path

$script:failures = 0
function Fail([string]$msg) {
  $script:failures++
  Write-Host "FAIL: $msg"
}
function Pass([string]$msg) { Write-Host "ok:   $msg" }

# Run the binary; returns @{ Code; Out } with stdout and stderr merged. Never throws on a non-zero exit.
function Akm {
  param([string]$Repo, [string[]]$Arguments)
  $prev = Get-Location
  Set-Location -LiteralPath $Repo
  $ErrorActionPreference = 'Continue'
  try {
    $out = & $Exe @Arguments 2>&1 | ForEach-Object { "$_" }
    $code = $LASTEXITCODE
  }
  finally { Set-Location $prev }
  return @{ Code = $code; Out = ($out -join "`n") }
}

# Expect an exit code and, optionally, a regex in the output.
function Expect {
  param([string]$Name, $Result, [int]$Code = 0, [string]$Match = '')
  if ($Result.Code -ne $Code) { Fail "$Name : exit $($Result.Code), want $Code`n$($Result.Out)"; return $false }
  if ($Match -and $Result.Out -notmatch $Match) { Fail "$Name : output does not match /$Match/`n$($Result.Out)"; return $false }
  Pass $Name
  return $true
}

function Invoke-Git([string]$Repo, [string[]]$Arguments) {
  $ErrorActionPreference = 'Continue'
  $out = & git -C $Repo @Arguments 2>&1 | ForEach-Object { "$_" }
  if ($LASTEXITCODE -ne 0) { throw "git $($Arguments -join ' ') failed: $($out -join "`n")" }
}

function WorkFile([string]$Repo, [string]$Slug, [string]$Sub = '') {
  $dir = Join-Path (Join-Path $Repo 'docs') 'work'
  if ($Sub) { $dir = Join-Path $dir $Sub }
  $f = Get-ChildItem -LiteralPath $dir -Filter "*-$Slug.md" -ErrorAction SilentlyContinue | Select-Object -First 1
  return $f
}

function Run-Mode([string]$AutoCrlf) {
  Write-Host "`n== core.autocrlf=$AutoCrlf =="
  # A path with spaces and nested directories, written with backslashes on Windows.
  $repo = Join-Path (Join-Path $Work "repo crlf-$AutoCrlf") 'sub dir'
  New-Item -ItemType Directory -Force -Path $repo | Out-Null

  Invoke-Git $repo @('init', '-q')
  Invoke-Git $repo @('config', 'core.autocrlf', $AutoCrlf)
  Invoke-Git $repo @('config', 'user.name', 'Smoke Test')
  Invoke-Git $repo @('config', 'user.email', 'smoke@example.invalid')

  $r = Akm $repo @('init', '--yes', '--no-input', '.')
  if (-not (Expect 'init' $r 0)) { return }
  foreach ($p in @('docs\.layout', 'docs\work\done', 'AGENTS.md')) {
    if (-not (Test-Path -LiteralPath (Join-Path $repo $p))) { Fail "init did not write $p" }
  }

  $r = Akm $repo @('new', 'work', 'smoke-alpha')
  if (-not (Expect 'new alpha' $r 0 'AKM-Item: (\d{6}-[0-9a-z]{4})')) { return }
  $alpha = [regex]::Match($r.Out, 'AKM-Item: (\d{6}-[0-9a-z]{4})').Groups[1].Value
  $r = Akm $repo @('new', 'work', 'smoke-beta')
  if (-not (Expect 'new beta' $r 0 'AKM-Item: (\d{6}-[0-9a-z]{4})')) { return }
  $beta = [regex]::Match($r.Out, 'AKM-Item: (\d{6}-[0-9a-z]{4})').Groups[1].Value

  # Put the files through git so core.autocrlf decides their line endings, as for a user who clones.
  Invoke-Git $repo @('add', '-A')
  Invoke-Git $repo @('commit', '-q', '-m', 'smoke: initial')
  Remove-Item -LiteralPath (Join-Path $repo 'docs') -Recurse -Force
  Invoke-Git $repo @('checkout', '--', '.')
  $alphaFile = WorkFile $repo 'smoke-alpha'
  $betaFile = WorkFile $repo 'smoke-beta'
  if (-not $alphaFile -or -not $betaFile) { Fail 'work items missing after checkout'; return }
  $bytes = [System.IO.File]::ReadAllBytes($alphaFile.FullName)
  $hasCrlf = ([System.Text.Encoding]::UTF8.GetString($bytes)) -match "`r`n"
  Write-Host "     record line endings after checkout: $(if ($hasCrlf) { 'CRLF' } else { 'LF' })"
  if ($AutoCrlf -eq 'true' -and $env:OS -eq 'Windows_NT' -and -not $hasCrlf) { Fail 'autocrlf=true did not produce CRLF; the case is not testing what it says' }

  # Fill, ready, claim.
  foreach ($id in @($alpha, $beta)) {
    $r = Akm $repo @('section', $id, 'Goal', '--text', 'Smoke goal.')
    Expect "section Goal $id" $r 0 | Out-Null
    $r = Akm $repo @('section', $id, 'Next', '--text', 'Smoke next.')
    Expect "section Next $id" $r 0 | Out-Null
    $r = Akm $repo @('set', $id, 'status=ready', 'priority=2')
    Expect "set ready $id" $r 0 "$id" | Out-Null
  }
  $r = Akm $repo @('claim', $alpha, '--as', 'smoke-agent')
  Expect 'claim alpha' $r 0 'active' | Out-Null
  $r = Akm $repo @('claim', $alpha, '--as', 'other-agent')
  Expect 'second claim is refused' $r 2 'claimed by smoke-agent' | Out-Null

  # A link from beta to alpha, with the record's own line ending.
  $alphaName = (WorkFile $repo 'smoke-alpha').Name
  $nl = if ($hasCrlf) { "`r`n" } else { "`n" }
  Add-Content -LiteralPath $betaFile.FullName -Value "${nl}See [alpha]($alphaName).$nl" -NoNewline

  # A repository path with mixed separators, and one given as a relative path with backslashes.
  $mixed = ($repo -replace '\\', '/') + '/./'
  $r = Akm $Work @('show', $alpha, '--repo', $mixed)
  Expect 'show with a mixed-separator --repo' $r 0 $alpha | Out-Null
  $r = Akm $Work @('ls', '--format', 'ids', '--repo', $repo)
  Expect 'ls with a native --repo' $r 0 $alpha | Out-Null

  # Path audit (work item 261108-kbve): the forms a Windows user types for --repo, and the paths printed back.
  $parent = Split-Path -Parent $repo
  $leaf = Split-Path -Leaf $repo
  $forms = [ordered]@{
    'trailing backslash'   = ($repo + '\')
    'forward slashes'      = ($repo -replace '\\', '/')
    'dot segments'         = ($repo + '\docs\..')
    'doubled separators'   = ($repo -replace '\\', '\\')
  }
  foreach ($k in $forms.Keys) {
    $r = Akm $Work @('show', $alpha, '--repo', $forms[$k])
    Expect "--repo with $k" $r 0 $alpha | Out-Null
  }
  # A relative --repo with a backslash, from the parent directory.
  $r = Akm $parent @('show', $alpha, '--repo', ".\$leaf")
  Expect '--repo relative with a backslash' $r 0 $alpha | Out-Null
  $r = Akm $repo @('find', $alpha)
  if ($r.Code -ne 0 -or $r.Out -match '\\') { Fail "find prints a path with a backslash or fails:`n$($r.Out)" } else { Pass 'find prints a forward-slash path' }
  # A note name written with a backslash is refused with a message, not a crash.
  $r = Akm $repo @('find', '--name', 'docs\work\smoke-alpha')
  Write-Host "     find --name with backslashes: exit $($r.Code): $($r.Out)"
  if ($r.Out -match '(?i)panic|goroutine') { Fail 'find --name with backslashes crashed' }

  # Close: moves the file into done/ and rewrites the link in beta.
  $r = Akm $repo @('section', $alpha, 'Outcome', '--text', 'Smoke outcome.')
  Expect 'section Outcome' $r 0 | Out-Null
  $r = Akm $repo @('close', $alpha, 'done')
  Expect 'close alpha' $r 0 'moved to' | Out-Null
  if (WorkFile $repo 'smoke-alpha') { Fail 'close left the record in docs/work' }
  if (-not (WorkFile $repo 'smoke-alpha' 'done')) { Fail 'close did not move the record to docs/work/done' }
  $betaText = [System.IO.File]::ReadAllText($betaFile.FullName)
  if ($betaText -notmatch 'done/' + [regex]::Escape($alphaName)) { Fail "close did not rewrite the link in beta:`n$betaText" } else { Pass 'close rewrote the link in beta' }
  if ($betaText -match '\\') { Fail 'a rewritten link holds a backslash' }

  # Record line endings stay consistent (no record mixes CRLF and LF).
  foreach ($f in Get-ChildItem -LiteralPath (Join-Path $repo 'docs\work') -Recurse -Filter '*.md') {
    $t = [System.IO.File]::ReadAllText($f.FullName)
    $crlf = ([regex]::Matches($t, "`r`n")).Count
    $lf = ([regex]::Matches($t, "(?<!`r)`n")).Count
    if ($crlf -gt 0 -and $lf -gt 0) { Fail "$($f.Name) mixes line endings ($crlf CRLF, $lf LF)" }
  }

  # A write while a reader holds the record open. Known fault 261108-mwtr: the rename onto an open
  # record is refused with "Access is denied". The harness asserts the safe outcome (the record is
  # intact, no temporary file is left, the next write works) and fails if the fault is fixed, so
  # the exception is removed with the fix.
  $r = Akm $repo @('set', $beta, 'priority=1')
  Expect 'set before the reader' $r 0 | Out-Null
  foreach ($case in @(
      @{ Name = 'a reader that shares read, write and delete'; Share = ([System.IO.FileShare]::ReadWrite -bor [System.IO.FileShare]::Delete) },
      @{ Name = 'a reader that shares read only'; Share = [System.IO.FileShare]::Read })) {
    $fs = [System.IO.File]::Open($betaFile.FullName, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, $case.Share)
    try { $r = Akm $repo @('set', $beta, 'priority=3') }
    finally { $fs.Dispose() }
    if ($r.Code -eq 0) {
      if ($env:OS -eq 'Windows_NT') { Fail "set with $($case.Name) succeeded: known fault 261108-mwtr is fixed, remove the exception" }
      else { Pass "set with $($case.Name)" }
    }
    elseif ($env:OS -eq 'Windows_NT') {
      if ($r.Out -match '(?i)panic|goroutine' -or $r.Out -notmatch 'akmctl: ') { Fail "set with $($case.Name) crashed or printed no error:`n$($r.Out)" }
      else { Pass "set with $($case.Name) was refused (known fault 261108-mwtr)" }
    }
    else { Fail "set with $($case.Name): exit $($r.Code)`n$($r.Out)" }
    if ([System.IO.File]::ReadAllText($betaFile.FullName) -notmatch '(?m)^priority: 1\r?$' -and $r.Code -ne 0) { Fail "a refused write changed the record ($($case.Name))" }
    if (Get-ChildItem -LiteralPath (Split-Path $betaFile.FullName) -Filter '.akm-*.tmp' -Force) { Fail "a temporary file is left after $($case.Name)" }
    $r = Akm $repo @('set', $beta, 'priority=1')
    Expect "set after $($case.Name) let go" $r 0 | Out-Null
  }

  # No lock file is left behind.
  if (Test-Path -LiteralPath (Join-Path $repo '.akm.lock')) { Fail 'a .akm.lock file is left behind' }

  $r = Akm $repo @('check')
  if ($AutoCrlf -eq 'true' -and $env:OS -eq 'Windows_NT') {
    # Known fault 261108-ma45: A4 compares the ADR index byte for byte, so a CRLF checkout fails.
    if ($r.Code -eq 0) { Fail 'check passed on the CRLF checkout: known fault 261108-ma45 is fixed, remove the exception' }
    elseif ($r.Out -match '1 errors' -and $r.Out -match 'A4 ') { Pass 'check reports only A4 on the CRLF checkout (known fault 261108-ma45)' }
    else { Fail "check on the CRLF checkout:`n$($r.Out)" }
  }
  else { Expect 'check' $r 0 '0 errors' | Out-Null }
  $r = Akm $repo @('ls', '--status', 'ready,active,blocked', '--format', 'ids')
  Expect 'ls open work' $r 0 $beta | Out-Null
}

try {
  $v = Akm $Work @('version')
  Expect 'version' $v 0 'akmctl' | Out-Null
  foreach ($mode in @('true', 'false')) { Run-Mode $mode }
}
finally {
  if ($env:AKM_SMOKE_KEEP -ne '1') { Remove-Item -LiteralPath $Work -Recurse -Force -ErrorAction SilentlyContinue }
}

if ($script:failures -gt 0) {
  Write-Host "`n$($script:failures) failure(s)"
  exit 1
}
Write-Host "`nwindows-smoke: all checks passed"
exit 0
