---
id: 261027-wkxv
type: task
status: done
origin: 261026-m74k
created: 2026-10-27
updated: 2026-10-28
authored_by: agent
title: Choose a YAML library for the Go tool
---

# Choose a YAML library for the Go tool

## Goal

The Go tool prefers the standard library, but a YAML library is needed to read frontmatter. Choose one by due diligence, in this order (decision of 2026-10-27):

1. **Baseline requirements.** Must parse the layout's frontmatter subset ([agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md), H5): flat keys, scalars, lists of scalars (inline and block style), integers, dates, quoted strings, long values. Must round-trip well enough for edits the tool makes, and report errors with line numbers.
2. **Performance.** Parse time and allocations on a few thousand small files.
3. **Security.** Known vulnerabilities, behaviour on hostile input (deep nesting, aliases and anchors, very large documents), maintenance record.
4. **Tiebreaker:** the more popular and better supported library.

Also check whether a small hand-written parser for the H5 subset would beat all of them, as `pkm-workbench` does today for flat keys.

Done when: a short note in the vault compares the candidates on those four points with evidence, names the choice, and says whether to adopt it or keep a hand-written subset parser.

## Outcome

Recommendation, awaiting Alexander's confirmation: keep a hand-written parser for the H5 subset (reads and edits by byte splice); no YAML library in `akmctl`. No library gave byte-for-byte edits (0 of 4), all accept out-of-subset input, and the hand parser was 6 times faster with bounded behaviour. If a library is ever wanted: `go.yaml.in/yaml/v3`. Use it in tests only as a differential oracle. Write-up, code and raw numbers: [go-yaml-library-evaluation](../../design/go-yaml-library-evaluation.md).

## Log

- 2026-10-27: filed by Claude Code on the Mac at Alexander's direction, split out of [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md).
- 2026-10-28: [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md) closed. The frontmatter subset it must parse is in section 0 of [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md) (H5), and [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) is accepted.
- 2026-10-28: claimed (Claude Code side session on the Mac, started by the AKM spec session; Build, scoped to this note and the design note).
- 2026-10-28: candidates `go.yaml.in/yaml/v3` v3.0.5, `/v4` v4.0.0-rc.6, `goccy/go-yaml` v1.19.2, `gopkg.in/yaml.v3` v3.0.1 (archived) and a hand-written parser tested for behaviour, edit round trip, speed and hostile input. Original `go-yaml/yaml` is archived; `yaml/go-yaml` is the maintained successor.
- 2026-10-28: decision (recommendation): hand-written subset parser; reason: edits byte-for-byte are impossible with every library and the subset check would be hand-written anyway. Hand parser bug found and fixed during testing (`k: - - x` accepted).
- 2026-10-28: ADR check: none (nothing changed yet; this is a recommendation). If Alexander confirms, the no-library choice may deserve a line in the tool's ADR; his call. [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) and the spec already allow a no-library tool (H5, X8). The spec could pin the subset details listed under open risks in [go-yaml-library-evaluation](../../design/go-yaml-library-evaluation.md).
- 2026-10-28: closed. Recommendation awaits Alexander's confirmation.
