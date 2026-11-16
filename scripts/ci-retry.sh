#!/bin/sh
# Retry helpers for CI steps that download from GitHub or a package feed (work item
# 261109-e540). Source this file: . scripts/ci-retry.sh
#
#   ci_retry <command...>      run the command, up to CI_RETRY_TRIES times (default 3) with
#                              a growing pause; a failure that is not a transient download
#                              error is not retried. Returns 0, the command's status, or 75
#                              when it still fails with a transient download error.
#   ci_skip_or_fail <status>   after a failed ci_retry: on 75 print a warning annotation and
#                              exit 0 (the step is skipped, the job goes on); otherwise exit
#                              with the status.
#
# Transient means the server or the network failed, not the change under test: an HTTP 500,
# 502, 503 or 504, a timeout, a reset connection. A checksum mismatch is never transient.

CI_RETRY_TRANSIENT='HTTP 50[0234]|error: 50[0234]|50[0234] \(|\(50[0234]\)|[Gg]ateway [Tt]ime-?out|Bad Gateway|Service Unavailable|timed out|Connection reset|connection was (reset|closed)|Could not resolve host|0x8019[0-9a-fA-F]{4}'

ci_is_transient() {
	grep -Eq "$CI_RETRY_TRANSIENT" "$1"
}

ci_retry() {
	ci_tries=${CI_RETRY_TRIES:-3}
	ci_pause=${CI_RETRY_PAUSE:-15}
	ci_log=$(mktemp)
	ci_n=1
	while :; do
		ci_status=0
		"$@" >"$ci_log" 2>&1 || ci_status=$?
		cat "$ci_log"
		if [ "$ci_status" -eq 0 ]; then
			rm -f "$ci_log"
			return 0
		fi
		if ! ci_is_transient "$ci_log"; then
			rm -f "$ci_log"
			return "$ci_status"
		fi
		if [ "$ci_n" -ge "$ci_tries" ]; then
			rm -f "$ci_log"
			return 75
		fi
		echo "ci-retry: try $ci_n of $ci_tries failed with a transient download error; waiting $((ci_pause * ci_n))s" >&2
		sleep $((ci_pause * ci_n))
		ci_n=$((ci_n + 1))
	done
}

ci_skip_or_fail() {
	if [ "$1" -eq 75 ]; then
		echo "::warning title=download failed::a download still failed after ${CI_RETRY_TRIES:-3} tries with a server or network error; this step is skipped and does not gate the run (work item 261109-e540)"
		exit 0
	fi
	exit "$1"
}
