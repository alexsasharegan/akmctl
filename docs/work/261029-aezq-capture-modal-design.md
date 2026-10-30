---
id: "261029-aezq"
title: capture modal design
summary: enhance form controls, parent picker, and style on the capture submit btn
status: raw
kind: feature
created: 2026-10-29
updated: 2026-10-29
---

# capture modal design

## Goal

1. modal header

current:

```
[new raw work item](small caps bold)
[Capture](heading)
```

Current text stacking does not fit normal hierarchy. Heading belongs on top of subtext/subheading.
Use the Radon font to exhibit the rawness of the capture in design as well as copy.
Place the subtitle below the heading. Consider muted look or italics to convey a secondary quality,
but prioritize accessibility legibility (contrast esp.).

proposed:

```
[[Raw](Radon) Capture](heading)
[new task record](small caps bold and muted)
```

2. modal size

modal size is constrained on desktop. could benefit from increased width. it's possible i just don't like how the form label descriptions wrap lines.

3. Summary field

160 characters does not fit nicely into this field and scrolls awkwardly.
consider a textedit field, or some way to see lines wrap.
consider adding helper text below input which counts characters and changes when input hits 160 or hits ~90% of max chars.

4. Kind

the input is not the same height as its siblings in the row.
plain html option/select inputs always feel like plain html inputs.
radix-ui or other popular libraries may serve as a good example of what custom UI selects feel like.

5. Parent

this is a massive opportunity for rich UI interactions, but also for big builds.
this needs a lot of runway to handle.
current search is squeezed into a small space,
does not offer any kind of emphasis on where matches are being made (in description or id),
and result options are too small to inspect to inform a good choice

## Log

- 2026-10-29: filed.
