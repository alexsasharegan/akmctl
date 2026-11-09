---
id: "261109-vxdr"
title: State what a MINOR and a MAJOR change mean
status: ready
priority: 1
created: 2026-11-09
updated: 2026-11-09
---

# State what a MINOR and a MAJOR change mean

## Goal

Exit criterion 5 of the 1.0.0 milestone says the specification states what a MINOR and a MAJOR change to the layout, the API and the extensions mean. No item covered it. Today the layout is 0.3.0, the API is 1.0 and the tool is 0.x, versioned independently (`docs/architecture.md`).

This needs the maintainer's decision, so the first step is a draft for them to accept, not a rule: write a `proposed` ADR with the options (for example: what counts as breaking for a reader of the records, for an HTTP or CLI client, and for an extension author; how long a MAJOR is supported; whether the three version numbers move together) and a recommendation.

Done when:

- [ ] A `proposed` ADR lists the options and consequences for the layout, the API and the extensions, with a recommendation.
- [ ] After the maintainer accepts it, the specification and the API spec carry the rule in one section each, and `docs/spec/akm-guide.md` points to it.
- [ ] Related: [261101-12c0](261101-12c0-standard-spec-enumerations.md) (standard spec enumerations) is read first and either folded in or left alone, with the reason in the Log.

## Next

Draft the proposed ADR for the maintainer to accept.

## Log

- 2026-11-09: filed.
- 2026-11-09 01:23-07:00: Filed: criterion 5 of the 1.0.0 milestone had no item.
