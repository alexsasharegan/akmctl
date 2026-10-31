---
id: "261031-st7s"
title: claim churn loop detection
summary: "Spike: detect an item that is claimed and released again and again without progress"
status: raw
kind: spike
parent: "261031-v07x"
created: 2026-10-31
updated: 2026-10-31
---

# claim churn loop detection

## Goal

From the consultation (gap 2). A loop of cold-start agents can flap: each turn rewrites `## Next` with a new guess, exits cleanly, and the item never closes. A counter in frontmatter would change a tracked file every turn. The claims extension already records `claimed_by` and `claimed_at`, and git history holds every claim and release. Done when:

- [ ] Compare: counting claim steps in git history (the step commits carry `AKM-Step` and `AKM-Item` trailers), a field in frontmatter, and leaving it to the harness. Show what each costs and what it can see (no Goal line ticked across N claims).
- [ ] A recommendation is in `## Outcome`, including whether a tripped item becomes `blocked` with `blocked_by` free text, and follow-up items are filed.

## Log

- 2026-10-31: filed from the consultation.
