---
id: 261028-h648
title: Offer a safety extension that scans records for known secret patterns and a declared deny list
summary: The optional safety extension (SF1 to SF4) reports secret-like strings and deny-list matches by file, line and rule, never the text. A pass is not proof.
status: accepted
date: 2026-10-28
scope: [internal/check]
authored_by: agent
---

# Offer a safety extension that scans records for known secret patterns and a declared deny list

## Context

Records can be copied from chat or logs and may carry credentials or private text. A complete check is impossible; a narrow one is cheap. (Decision log entries Y13, Y24, Y36.)

## Considered Options

- No check: leaks are caught by review only.
- Known patterns plus a repository-declared deny list, with reasoned allow markers.
- Entropy scanning alone: noisy.

## Decision

Rules SF1 (known patterns), SF2 (`docs/.akm-deny`, literals or `re:` lines), SF3 (high-entropy strings, skipping hashes and URLs) and SF4 (inline `akm-allow: SFn <reason>` markers; one without a reason is an error). Findings never print the matched text.

## Consequences

- Known formats only; a pass means no known pattern.
- This repository enables the extension and scans its own records.

## Confirmation

`internal/check/safety_test.go`; `akmctl check` on this repository.
