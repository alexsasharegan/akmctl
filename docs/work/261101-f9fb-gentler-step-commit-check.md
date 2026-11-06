---
id: "261101-f9fb"
title: gentler step commit check
status: raw
created: 2026-11-01
updated: 2026-11-06
---

# gentler step commit check

## Goal

A maintainer commit with the subject `261101-4gk3(file): ...` and no AKM- trailers failed CM2 ("found 0 and 0"), and only a history rewrite could clear it. Codex hit a related case (261101-ft79). Make `commits` strict about what a commit claims and lenient about whether it claims anything, so a person who types the natural subject is not punished. Done when:

- [ ] A commit whose subject is step-shaped but has no AKM- trailers at all is reported as a warning (and a CM4-style check still flags a status or claim change), not as a CM2 failure; the exit code stays 0 for it.
- [ ] A commit with a wrong, duplicated or partial trailer set (one of the pair, or a mismatch with the subject) still fails CM2.
- [ ] The spec CM table, the `commits` skill reference and the CHANGELOG say so, with tests; fold in 261101-ft79 (trailers outside git's trailer block) or link it.
- [ ] Decide whether this needs an ADR: it changes a check that agents rely on.

## Log

- 2026-11-01: filed.
- 2026-11-06 13:19-07:00: Detached from 261101-vs5c on 2026-11-06: it adds a warning severity to commits, an API change that probably needs an ADR, so it is its own item, not a rails gap.
