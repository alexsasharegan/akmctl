---
id: 261028-fnez
title: Specify the operations once and bind them to the CLI and to loopback HTTP with a bearer token
summary: One transport-neutral API with schemas and fixtures; the CLI and an HTTP server on 127.0.0.1 with a per-run token implement it. Read before adding an operation.
status: accepted
date: 2026-10-28
scope: [internal/api, internal/server, internal/cli, testdata/api]
authored_by: agent
---

# Specify the operations once and bind them to the CLI and to loopback HTTP with a bearer token

## Context

Other implementations and the UI must be able to rely on the same operations and the same JSON. The UI needs writes, not only reads. (Decision log entry Y34.)

## Considered Options

- Read-only HTTP: the UI would shell out to the CLI for every change.
- OpenAPI: more tooling, and the CLI binding needs its own description.
- Prose plus JSON Schema 2020-12, one schema per operation, with conformance fixtures; HTTP with writes on loopback only.

## Decision

The API spec is prose plus JSON Schema with fixtures that validate against it. The CLI and the HTTP server are two bindings of `internal/api`. HTTP binds 127.0.0.1, takes a per-run bearer token, checks Host and Origin, sends no CORS headers, and writes carry an optional `if_match` etag. The fixtures in `testdata/api/` are a copy of the spec's; the spec copy is the source.

## Consequences

- Any operation is reachable the same way from both bindings.
- The conformance test in `internal/api` runs 25 cases against sample repositories.
- The fixtures exist twice and are copied by hand in both directions.

## Confirmation

`internal/api` conformance test.
