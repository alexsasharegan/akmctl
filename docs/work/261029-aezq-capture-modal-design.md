---
id: "261029-aezq"
title: "Capture modal: header, width, a cycle field and the submit button"
summary: Rework the capture modal's header and width, add a Cycle field, and restyle the submit button; select, summary and parent controls come from separate items
status: ready
kind: feature
created: 2026-10-29
updated: 2026-10-29
priority: 1
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

- [ ] the header reads heading first, with the muted subtitle below, and the heading uses Radon; contrast passes in light and dark
- [ ] the modal is wider on desktop and its label descriptions no longer wrap awkwardly
- [ ] an optional Cycle field files the item with that cycle, and a bad label shows the API's error on the field
- [ ] the submit button has its own style with visible focus in both themes
- [ ] unit tests cover the cycle field's request; checked in the browser in light and dark
- [ ] new terms are registered in docs/glossary.md

## Next

Ready. Do 261029-kb33 first if the select and summary should look right inside the new layout.

## Log

- 2026-10-29: filed.
- 2026-10-29: refined with Alexander; the summary, kind and parent parts moved to separate items.
