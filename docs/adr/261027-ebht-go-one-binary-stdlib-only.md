---
id: 261027-ebht
title: Write the tool in Go as one binary using only the standard library
summary: The tool is Go, one portable binary, standard library only at run time; a YAML library is allowed in tests as an oracle. Read before adding a dependency.
status: accepted
date: 2026-10-27
scope: [go.mod, cmd/akmctl]
authored_by: agent
---

# Write the tool in Go as one binary using only the standard library

## Context

The standard needs a tool that queries and checks records and also serves a workbench UI. The first plan was Python. The maintainer already works in Go, wants a single prebuilt binary with no runtime on the user's machine, and wants the CLI and the UI in one codebase. (Decision log entries X8, Y9, Y28.)

## Considered Options

- Python: quick to write, needs an interpreter and packages on every machine.
- Go with a YAML library and a width library: less code, two dependencies to vet and update.
- Go, standard library only: more code to write, nothing to vet at run time.
- Do nothing: no tool; the standard would stay a text only.

## Decision

Go, one module, standard library only in the binary. `go.yaml.in/yaml/v3` is used in tests only, as the oracle for the parser (ADR 261028-h99k). Display width is the tool's own code (ADR 261028-mw1r).

## Consequences

- One binary per platform; no runtime to install; cross-compiles for five targets.
- The parser, the width rule and the HTTP binding are ours to maintain.
- A dependency needs an ADR that says why the standard library does not do.

## Confirmation

`go.mod` lists only `go.yaml.in/yaml/v3`, and `go list -deps ./cmd/akmctl` shows no third-party package.
