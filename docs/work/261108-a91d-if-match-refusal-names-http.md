---
id: "261108-a91d"
title: if match refusal names http
status: raw
created: 2026-11-08
updated: 2026-11-08
---

# if match refusal names http

## Goal

A write over HTTP without `If-Match` is refused with status 428, but the response says `"operation": "http"` and `"request": {}` (the route layer fails before the operation runs, `internal/server/server.go` `route` then `fail`), so a client cannot tell which operation was refused, and the other refusals of a request name the operation and echo the request. Section 6 of the API spec does not say what the envelope of this refusal holds. Found by 261108-f1tv; the fixture `set-15` records today's answer.

Done when:

- [ ] The spec says what `operation` and `request` hold in the `if_match_required` refusal (and in the other route-level failures).
- [ ] The server answers as the spec says, and `set-15` and the schemas follow.

## Next

Decide whether the refusal should name the operation and echo the request; propose it to the maintainer before changing the server or the fixture.

## Log

- 2026-11-08: filed.
