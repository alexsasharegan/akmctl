---
id: "261029-ap3y"
title: Keep docs/.akm.lock out of commits
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: chore
priority: 3
summary: The lock file got committed by a broad git add in the vault run.
---

# Keep docs/.akm.lock out of commits

## Goal

akmctl writes docs/.akm.lock, which the layout says is not committed. Write the ignore line when the tool first creates it, or document it in the guide.

Done when:

- [ ] A repository's .gitignore lists the lock file after the first write.
- [ ] The guide says so.

## Log

- 2026-10-29: filed.
