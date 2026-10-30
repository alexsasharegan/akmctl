---
id: "261029-7zcp"
title: Unblocking a claimed item to active logs a second claim
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
---

# Unblocking a claimed item to active logs a second claim

## Goal

Moving a blocked item that already has a claim to active makes the claim op run again: the log gets a second 'claimed' line. Skip the claim when the same identity already holds it; ask for a takeover when another holds it. Found in the dogfood pass.

## Next

Review the claim change in internal/api/ops_write.go; then close with the drafted Outcome.

## Outcome

Re-claiming an active item you already hold refreshes `claimed_at` and adds no log line; another holder still needs `takeover`. Test in `internal/api/write_test.go`.

## Log

- 2026-10-29: filed.
- 2026-10-29: fixed in claim: the same holder re-claiming an active item refreshes claimed_at and logs nothing; another holder still needs takeover. Test TestUnblockClaimedItemLogsOneClaim. Item is raw so it was not claimed.
