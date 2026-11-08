---
id: "261108-sfj3"
title: backslash link targets are ignored
status: raw
parent: 261108-2ncf
created: 2026-11-08
updated: 2026-11-08
---

# backslash link targets are ignored

## Goal

Found by the CRLF and path audit ([261108-kbve](done/261108-kbve-windows-crlf-and-path-audit.md)), part of [261108-2ncf](261108-2ncf-basic-windows-support.md).

Input: a record that links with Windows separators, such as `[x](..\adr\README.md)` or `[y](done\261108-aaaa-slug.md)`.

Actual: `check` reports nothing for it and `close` does not rewrite it (it stays as typed, so it breaks after a move). Forward-slash links are checked and rewritten.

Expected: decide the rule. Either `check` reports a backslash in a link target (rule and message) or the tool reads it as a path and rewrites it with forward slashes; the standard says records use forward slashes.

Done when:

- [ ] The rule is decided and tested on any platform (a backslash is a plain character off Windows).
- [ ] The guide says what a Windows user should type.

## Next

Probe check and close with a backslash link on a copy of the fixture; decide report versus rewrite. Low priority.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:52-07:00: filed by the CRLF and path audit
