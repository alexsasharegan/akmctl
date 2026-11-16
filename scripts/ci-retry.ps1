# Retry helpers for CI steps that download from GitHub or the winget feed (work item
# 261109-e540). Dot-source this file: . ./scripts/ci-retry.ps1
#
#   Invoke-CiRetry { <command> }   run the script block, up to 3 times with a growing pause;
#                                  a failure that is not a transient download error is not
#                                  retried. Returns @{ Ok; Transient; Output; ExitCode }.
#   Exit-CiSkipOrFail $result      after a failed Invoke-CiRetry: a transient failure prints a
#                                  warning annotation and exits 0 (the step is skipped, the job
#                                  goes on); any other failure exits with its code.
#
# Transient means the server or the network failed, not the change under test: an HTTP 500,
# 502, 503 or 504, a timeout, a reset connection. A checksum mismatch is never transient.

$CiRetryTransient = 'HTTP 50[0234]|error: 50[0234]|50[0234] \(|\(50[0234]\)|Gateway Time-?out|Bad Gateway|Service Unavailable|timed out|Connection reset|connection was (reset|closed)|could not be resolved|0x8019[0-9a-fA-F]{4}'

function Test-CiTransient([string]$Text) {
    return $Text -match $CiRetryTransient
}

function Invoke-CiRetry([scriptblock]$Command, [int]$Tries = 3, [int]$Pause = 15) {
    $out = ''
    $code = 1
    foreach ($n in 1..$Tries) {
        $global:LASTEXITCODE = 0
        $out = & $Command 2>&1 | Out-String
        $code = $LASTEXITCODE
        Write-Host $out
        if ($code -eq 0) { return @{ Ok = $true; Transient = $false; Output = $out; ExitCode = 0 } }
        if (-not (Test-CiTransient $out)) { return @{ Ok = $false; Transient = $false; Output = $out; ExitCode = $code } }
        if ($n -lt $Tries) {
            Write-Host "ci-retry: try $n of $Tries failed with a transient download error; waiting $($Pause * $n)s"
            Start-Sleep -Seconds ($Pause * $n)
        }
    }
    return @{ Ok = $false; Transient = $true; Output = $out; ExitCode = $code }
}

function Exit-CiSkipOrFail($Result) {
    if ($Result.Transient) {
        Write-Host "::warning title=download failed::a download still failed after the retries with a server or network error; this step is skipped and does not gate the run (work item 261109-e540)"
        exit 0
    }
    exit $Result.ExitCode
}
