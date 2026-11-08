---
id: "261029-aezq"
title: "Capture modal: header, width, a cycle field and the submit button"
summary: Rework the capture modal's header and width, add a Cycle field, and restyle the submit button; select, summary and parent controls come from separate items
status: active
kind: feature
created: 2026-10-29
updated: 2026-11-07
priority: 1
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T05:34
---

# Capture modal: header, width, a cycle field and the submit button

## Goal

Polish the capture modal (key `c`). The styling freeze is lifted for this item; use the existing design tokens. The select and summary controls are in 261029-kb33, and the parent picker is in 261029-bdt9.

1. Modal header

Current:

```
[new raw work item](small caps bold)
[Capture](heading)
```

Current text stacking does not fit normal hierarchy. Heading belongs on top of subtext/subheading. Use the Radon font to exhibit the rawness of the capture in design as well as copy. Place the subtitle below the heading. Consider a muted look or italics to convey a secondary quality, but prioritize accessible legibility (contrast especially).

Proposed:

```
[[Raw](Radon) Capture](heading)
[new task record](small caps bold and muted)
```

2. Modal size

The modal size is constrained on desktop and could benefit from increased width. It is possible the real complaint is that the form label descriptions wrap lines.

3. Cycle field (added in refinement)

The form offers Title, Name, Kind, Priority, Parent and Summary today. Add an optional Cycle label next to Priority, checked against the same pattern `set` uses, so a team can file straight into a sprint.

4. Submit button

The capture submit button needs a style of its own (named in the original note; details to agree when this is picked up).

Moved to other items: the summary field and the kind select (261029-kb33); the parent search (261029-bdt9).

Done when:

- [x] the header reads heading first, with the muted subtitle below, and the heading uses Radon; contrast passes in light and dark
- [x] the modal is wider on desktop and its label descriptions no longer wrap awkwardly
- [x] an optional Cycle field files the item with that cycle, and a bad label shows the API's error on the field
- [x] the submit button has its own style with visible focus in both themes
- [x] unit tests cover the cycle field's request; checked in the browser in light and dark
- [x] new terms are registered in docs/glossary.md

## Next

All boxes are ticked: close it once the pull request is merged or the maintainer is happy with the look. The dark-mode focus ring on the submit button was not checked by eye (it uses the `--focus` token).

## Log

- 2026-10-29: filed.
- 2026-10-29: refined with Alexander; the summary, kind and parent parts moved to separate items.
- 2026-11-07 22:34-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 22:37-07:00: ticked "the header reads heading first, with the muted subtitle below, and the heading uses Radon; contrast passes in light and dark": heading first with the muted New work item line below; Raw in Radon; contrast 9.8 and 7.8 to 1 in dark and light
- 2026-11-07 22:37-07:00: ticked "the modal is wider on desktop and its label descriptions no longer wrap awkwardly": dialog 720px (was 560); parent search on its own row; checked at 375px
- 2026-11-07 22:37-07:00: ticked "an optional Cycle field files the item with that cycle, and a bad label shows the API's error on the field": new now accepts fields.cycle; the modal filed cycle: 2026-s3; a bad label shows the error on the field
- 2026-11-07 22:37-07:00: ticked "the submit button has its own style with visible focus in both themes": cap-go style; 3px focus ring checked in light (dark uses --focus), button contrast 8.8 to 1
- 2026-11-07 22:37-07:00: ticked "unit tests cover the cycle field's request; checked in the browser in light and dark": buildCapture cycle cases in logic.test.mjs, TestNewWorkWithCycle; browser checked in light and dark
- 2026-11-07 22:37-07:00: ticked "new terms are registered in docs/glossary.md": glossary Capture modal row extended
- 2026-11-07: reworked the modal: heading first with Raw in Radon and the muted subtitle, 720px wide, Cycle field beside Priority, Parent on its own row, and a `cap-go` submit style. `new` now accepts `fields.cycle` for a work item (spec, test, changelog), which the Cycle field needs.
