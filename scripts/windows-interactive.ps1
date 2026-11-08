<#
.SYNOPSIS
  What the interactive commands do where there is no console: init without a terminal, and akmctl ui.

.DESCRIPTION
  The GitHub Windows runner has no console a person can type into, so this cannot run the
  init questionnaire the way Windows Terminal does. It checks the cases a runner can:
  init with stdin from NUL and from a pipe (neither may hang or fail), and akmctl ui
  (the page is served on the loopback interface and answers over HTTP; with --open the
  server still starts). Work item 261108-ct4s. What it cannot do is listed in
  docs/windows-first-run.md and docs/architecture.md.

.EXAMPLE
  pwsh -File scripts\windows-interactive.ps1 -Exe $env:TEMP\akmctl.exe
#>
param([Parameter(Mandatory = $true)][string]$Exe)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$Exe = (Resolve-Path -LiteralPath $Exe).Path
$Work = Join-Path ([System.IO.Path]::GetTempPath()) ("akm interactive " + [guid]::NewGuid().ToString('N').Substring(0, 8))
New-Item -ItemType Directory -Force -Path $Work | Out-Null
$script:failures = 0
function Fail([string]$m) { $script:failures++; Write-Host "FAIL: $m" }
function Pass([string]$m) { Write-Host "ok:   $m" }

# init with stdin redirected; the process must finish by itself in 60 seconds.
function Init-With([string]$Name, [string]$Stdin) {
  $dir = Join-Path $Work $Name
  New-Item -ItemType Directory -Force -Path $dir | Out-Null
  $outFile = Join-Path $Work "$Name.out"
  $errFile = Join-Path $Work "$Name.err"
  # cmd.exe does the redirection, as a user or a script would.
  $cmdline = '"' + $Exe + '" init "' + $dir + '" < "' + $Stdin + '" > "' + $outFile + '" 2> "' + $errFile + '"'
  $p = Start-Process -FilePath cmd.exe -ArgumentList @('/c', ('"' + $cmdline + '"')) -PassThru -NoNewWindow
  if (-not $p.WaitForExit(60000)) { & taskkill /F /T /PID $p.Id | Out-Null; Fail "init with $Name did not finish in 60 seconds (it waits for input)"; return }
  $text = (Get-Content -Raw $outFile) + (Get-Content -Raw $errFile)
  Write-Host "     init with ${Name}: exit $($p.ExitCode), asked: $($text -match '\[[Yy]/[Nn]\]')"
  if ($p.ExitCode -ne 0) { Fail "init with $Name exit $($p.ExitCode)`n$text"; return }
  if (-not (Test-Path -LiteralPath (Join-Path $dir 'AGENTS.md'))) { Fail "init with $Name wrote no AGENTS.md"; return }
  Pass "init with $Name finished and wrote the repository"
}

# stdin from NUL (a character device on Windows, as /dev/null is on Unix).
Init-With 'stdin NUL' 'NUL'
# stdin from an empty file (a regular file, like a script that closes its input).
$empty = Join-Path $Work 'empty.txt'
New-Item -ItemType File -Force -Path $empty | Out-Null
Init-With 'stdin empty file' $empty
# stdin from a file of answers: not a terminal, so nothing is asked and the answers are not read.
$ans = Join-Path $Work 'answers.txt'
"n`nn`nn`n" | Set-Content $ans
Init-With 'stdin answers file' $ans

# akmctl ui: serve, fetch, stop. Called for plain and --open.
function Ui-Case([string]$Name, [string[]]$Extra) {
  $repo = Join-Path $Work "ui $Name"
  New-Item -ItemType Directory -Force -Path $repo | Out-Null
  & $Exe init --yes --no-input $repo | Out-Null
  $errFile = Join-Path $Work "ui-$Name.err"
  $outFile = Join-Path $Work "ui-$Name.out"
  $uiArgs = @('ui', '--repo', ('"' + $repo + '"')) + $Extra
  $p = Start-Process -FilePath $Exe -ArgumentList $uiArgs -RedirectStandardError $errFile -RedirectStandardOutput $outFile -PassThru -NoNewWindow
  try {
    $url = $null
    for ($i = 0; $i -lt 100 -and -not $url; $i++) {
      Start-Sleep -Milliseconds 100
      if (Test-Path $errFile) {
        $m = [regex]::Match((Get-Content -Raw $errFile), 'open: (http://[^\s]+)')
        if ($m.Success) { $url = $m.Groups[1].Value }
      }
      if ($p.HasExited) { break }
    }
    if (-not $url) { Fail "ui $Name printed no URL (exit $($p.HasExited)):`n$(Get-Content -Raw $errFile)"; return }
    Pass "ui $Name printed a loopback URL"
    if ($url -notmatch '^http://127\.0\.0\.1:\d+/#token=') { Fail "ui $Name URL is not on 127.0.0.1 with a token: $url" }
    $page = $url -replace '#.*$', ''
    $r = Invoke-WebRequest -UseBasicParsing -Uri $page -TimeoutSec 20
    if ($r.StatusCode -ne 200 -or $r.Content -notmatch '(?i)<html') { Fail "ui $Name page: status $($r.StatusCode)" }
    else { Pass "ui $Name served the page ($($r.Content.Length) bytes)" }
    # The API refuses a request without the token.
    try {
      $a = Invoke-WebRequest -UseBasicParsing -Uri ($page.TrimEnd('/') + '/api/v1/version') -Method Post -Body '{}' -ContentType 'application/json' -TimeoutSec 20
      Write-Host "     api without a token: status $($a.StatusCode)"
      if ($a.StatusCode -eq 200) { Fail "ui $Name answered an API call without the token" }
    }
    catch {
      $code = [int]$_.Exception.Response.StatusCode
      Write-Host "     api without a token: status $code"
      if ($code -eq 200) { Fail "ui $Name answered an API call without the token" } else { Pass "ui $Name refused an API call without the token ($code)" }
    }
    if ($p.HasExited) { Fail "ui $Name exited early with $($p.ExitCode)" }
  }
  finally {
    if (-not $p.HasExited) { $p.Kill() }
    $p.WaitForExit(10000) | Out-Null
  }
}
Ui-Case 'plain' @()
Ui-Case 'open' @('--open')

Remove-Item -LiteralPath $Work -Recurse -Force -ErrorAction SilentlyContinue
if ($script:failures -gt 0) { Write-Host "`n$($script:failures) failure(s)"; exit 1 }
Write-Host "`nwindows-interactive: all checks passed"
exit 0
