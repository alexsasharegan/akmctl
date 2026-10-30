---
id: "261029-4hx0"
title: Add an SVG favicon
status: done
created: 2026-10-29
updated: 2026-10-30
kind: feature
priority: 3
summary: The workbench has no favicon, so the browser tab shows a blank icon.
---

# Add an SVG favicon

## Goal

Add a favicon as an SVG served from the embedded assets (the asset server serves only listed prefixes and extensions, so add a path such as /icons/favicon.svg, or allow .svg under an existing prefix, with the right content type) and link it from index.html. Design: a simple mark that fits the Quiet paper style and reads at 16 px; it should work in light and dark (use prefers-color-scheme inside the SVG, or a colour that works on both). Optional, noted for later: a per-repository colour accent so two tabs differ at a glance, derived from a hash of the repository name.

Done when:

- [ ] favicon.svg is embedded, served without a token, with the image/svg+xml content type, and linked in the page head.
- [ ] The CSP still has no unsafe-inline (an SVG in img or link is fine).
- [ ] It reads at 16 px in light and dark browser chrome.
- [ ] A server test covers the route and its content type.

## Next

- Coordinator: look at the tab icon at 16 px in light and dark browser chrome (not done here; no browser was run). Then tick the last Done item.
- The per-repository colour accent (optional, in the Goal) is not done.

## Outcome

`internal/server/web/favicon.svg` (light and dark) is linked from `index.html` and served at `/favicon.svg` as image/svg+xml; `TestFaviconIsServed` covers it. Link resolved in the built-in browser; the 16px rendering was not inspected.

## Log

- 2026-10-29: filed.
- 2026-10-30: claimed (ui-agent).
- 2026-10-30: web/favicon.svg is a rounded tile with three bars, no text. Its colours switch with prefers-color-scheme inside the SVG. Served at `/favicon.svg` by an explicit case in server.go (image/svg+xml, same CSP, no token) and linked from index.html. TestFaviconIsServed covers the route, the type and the link.
- 2026-10-30: released (ui-agent).
