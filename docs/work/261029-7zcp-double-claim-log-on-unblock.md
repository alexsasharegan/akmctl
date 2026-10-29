---
id: "261029-7zcp"
title: Unblocking a claimed item to active logs a second claim
status: raw
created: 2026-10-29
updated: 2026-10-29
---

# Unblocking a claimed item to active logs a second claim

## Goal

Moving a blocked item that already has a claim to active makes the claim op run again: the log gets a second 'claimed' line. Skip the claim when the same identity already holds it; ask for a takeover when another holds it. Found in the dogfood pass.

## Log

- 2026-10-29: filed.
