---
id: 261028-h99k
title: Parse and edit frontmatter with a hand-written YAML-lite parser, not a YAML library
summary: Frontmatter is YAML-lite, a strict subset of YAML 1.2 core, read and edited by our parser that splices bytes. Read before touching internal/frontmatter.
status: accepted
date: 2026-10-28
scope: [internal/frontmatter, testdata/yaml-lite.json]
authored_by: agent
---

# Parse and edit frontmatter with a hand-written YAML-lite parser, not a YAML library

## Context

Records carry frontmatter that the tool must read exactly and edit without disturbing the rest of the file: comments, quoting, key order, line endings. Obsidian reads and rewrites the same files, and its YAML is 1.2 core. A YAML library decodes to a value and re-encodes the whole block, losing comments and formatting. (Decision log entries Y8, Y9, Y15, Y26.)

## Considered Options

- A YAML library (`go.yaml.in/yaml/v3`) with a node API: heavy, and rewriting loses bytes.
- A hand-written parser for a strict subset: small, exact errors with line numbers, byte-splice edits.
- Do nothing: no parsing, no tool.

## Decision

A hand-written parser for the subset in the specification (flat keys, scalars, lists of scalars; errors for everything richer) and an editor that changes only the bytes of the entry it touches. Dates, `yes`, `no`, `on`, `off` and `12:30` are plain strings, as in YAML 1.2 core and as Obsidian writes them. The YAML library runs in tests only, as an oracle over 52 fixtures plus a fuzz test.

## Consequences

- Every unsupported YAML feature is an error, so a record is never read differently by another parser.
- A library bump cannot change what the tool accepts.
- The fixtures in `testdata/yaml-lite.json` are shared with the specification and bind every implementation.

## Confirmation

`internal/frontmatter` tests: the 52 fixtures, the oracle test and the fuzz test.
