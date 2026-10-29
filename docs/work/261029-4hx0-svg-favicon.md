---
id: "261029-4hx0"
title: Add an SVG favicon
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: The workbench has no favicon, so the browser tab shows a blank icon.
---

# Add an SVG favicon

## Goal

Add a favicon as an SVG served from the embedded assets (a path under the allowed asset prefixes, for example /fonts-free /icons/favicon.svg or a new allowed prefix; the asset server allows only listed prefixes and extensions, so add .svg to the allowlist with the right content type) and link it from index.html. Design: a simple mark that fits the Quiet paper style and reads at 16 px; it should work in light and dark (use prefers-color-scheme inside the SVG, or a colour that works on both). Optional, noted for later: a per-repository colour accent so two tabs differ at a glance, derived from a hash of the repository name.

Done when:

- [ ] favicon.svg is embedded, served without a token, with the image/svg+xml content type, and linked in the page head.
- [ ] The CSP still has no unsafe-inline (an SVG in img or link is fine).
- [ ] It reads at 16 px in light and dark browser chrome.
- [ ] A server test covers the route and its content type.

## Log

- 2026-10-29: filed.
