---
id: "261029-g8x2"
title: The review diff is against the current file, not after earlier writes
summary: Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-09
closed: 2026-11-09
---

# The review diff is against the current file, not after earlier writes

## Goal

Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.

Done when:

- [x] The review shows each diff as it would apply after the earlier writes.

## Next

## Outcome

The move wizard's review shows each diff as it will apply after the earlier writes. A new read-only API operation, preview (ADR 261109-vgg7, proposed), applies the writes on a scratch copy; the wizard makes one call to it and previews every write.

## Log

- 2026-10-29: filed.
- 2026-11-09 10:37-07:00: claimed (claude-sonnet-5-5-g8x2).
- 2026-11-09 10:42-07:00: ticked "The review shows each diff as it would apply after the earlier writes.": TestPreviewShowsEachWriteAfterTheEarlierOnes; wizard.test.mjs and keys.test.mjs; conformance cases preview-1 to preview-3
- 2026-11-09 10:42-07:00: Chose a server preview operation (ADR 261109-vgg7, proposed): POST /preview applies a write list on a scratch copy and returns each write's dry-run response. The wizard review makes one call; previewable is gone, so a claim after a status change and a set naming the new ADR now have diffs.
