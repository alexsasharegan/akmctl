---
id: "261029-mj49"
title: After a capture, offer to open the new record
status: raw
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 0
summary: "After the capture modal files a raw record, a toast offers to open it: an Open action, or a temporary key binding the toast names"
---

# After a capture, offer to open the new record

## Goal

After a raw record is captured in the workbench, offer a follow-up open of the new file. Alexander's two options: a toast with an action button once toast actions exist, or a toast that shows a temporary key binding (for example `o` for the length of the toast).

Open questions:

- Open where: the drawer (the record's detail view), or Obsidian through the existing `obsidianUrl` helper (logic.js), or both.
- Depends on 261029-kmre for the toast shell; the key-binding variant needs no action buttons, so it could ship first.
- Whether `akmctl new` on the CLI should print an open hint too (probably not).

Done when:

- [ ] after a capture the toast offers to open the new record, by an action or a named key
- [ ] the offer expires with the toast and never steals a key that has another meaning at that moment
- [ ] unit tests cover the key's lifetime; checked in the browser
- [ ] new terms are registered in docs/glossary.md

## Log

- 2026-10-29: filed.
