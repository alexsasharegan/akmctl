---
id: "261108-jbd5"
title: remove has no http route
status: raw
created: 2026-11-08
updated: 2026-11-08
---

# remove has no http route

## Goal

The operation `remove` (API 1.0) has no HTTP route: the route table in section 6 of the API spec and `internal/server/server.go` have no entry for it, so the conformance runner skips the six `remove` cases on the HTTP binding (found by 261108-f1tv). Every other operation has a route.

Done when:

- [ ] The maintainer decides: add `DELETE /records/{id}/sections/{name}/items/{n}` (or the route that fits) to the spec and the server, or state in section 6 that `remove` is API and CLI only.
- [ ] If a route is added: the six `remove` cases run on the HTTP binding in `internal/conformance` and pass; if not, the runner's skip reason points at the spec sentence.

## Next

Read section 4.23 of the API spec for `remove` and propose the route; ask the maintainer before adding it to the spec.

## Log

- 2026-11-08: filed.
