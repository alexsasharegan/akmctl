---
id: "261030-43jw"
title: Evaluate and potentially migrate to Oxc toolchain (oxlint / oxfmt)
status: ready
created: 2026-10-30
updated: 2026-10-30
kind: spike
summary: Compare oxlint and oxfmt with the current Prettier setup for compatibility and speed, then decide whether to migrate
---

# Evaluate and potentially migrate to Oxc toolchain (oxlint / oxfmt)

## Goal

### Context

A migration to the Rust-based Oxc toolchain (`oxlint` and `oxfmt`) has been proposed to replace the current ESLint and Prettier setup. The tools are built for the "Vite stack" by VoidZero and compile to native binaries, promising up to 20x faster execution and near-instant cold starts.

This item evaluates their compatibility with `akmctl`'s current rule sets and measures the actual performance gain before committing to a full replacement.

### Objectives

- Determine whether `oxfmt` can act as a 1:1 drop-in replacement for Prettier.
- Determine whether `oxlint` covers the necessary ESLint rules, including TypeScript-specific rules.
- (Optional) Benchmark the difference in local execution and pre-commit hooks.

### Done when

- [ ] **Audit:** the Log records the current setup. Prettier 3.9.10 is the only formatter or linter, and it checks Markdown only (`AGENTS.md`, `CHANGELOG.md`, `README.md`, `docs/**`). There is no ESLint, and the steps below that assume husky or lint-staged are re-scoped or dropped.
- [ ] **Format comparison:** `oxfmt` ran over that same file set, and the Log records its diff against `prettier --check` (the drift, or "none"). It also records whether `oxfmt` formats Markdown at all.
- [ ] **Lint comparison:** `oxlint` ran over the repository's JavaScript (the workbench UI and `docs/spec/akm-api/validate.js`), and the Log lists the rules it covers and any it lacks.
- [ ] **Decision:** the Outcome records migrate, partial or stay, with the evidence. A migrate or partial decision files follow-up work items, and an ADR draft (`proposed`) is written if Prettier is rejected.
- [ ] **Benchmark (optional):** cold and cached timings for both setups are in the Log.

### Implementation steps (if approved)

- [ ] Remove the ESLint and Prettier dependencies (`eslint`, `prettier` and related plugins).
- [ ] Add `oxlint` and `oxfmt` to the `package.json` scripts.
- [ ] Update the `.husky` / `lint-staged` pre-commit hooks to use the new binaries.
- [ ] Update CI (for example GitHub Actions) to use the Oxc tools.
- [ ] Update `CONTRIBUTING.md` or the equivalent developer documentation.

### Notes and considerations

- **Type-aware rules:** `oxlint` parses natively, but if we rely heavily on ESLint's `type-checked` rules we may need to run `tsc --noEmit` alongside it.
- **Prettier plugins:** if non-standard files are formatted with Node-based Prettier plugins, `oxfmt` may not support them natively.
- **Fallback:** if the two are a clean switch, keep Prettier and ESLint as fallbacks for when the Vite stack is unavailable to a user (for example company-quarantined mirrors).

## Log

- 2026-10-30: filed.
- 2026-10-30: scope note for the audit: `package.json` has only `prettier` (checks Markdown records) and `ajv`; no ESLint, husky or lint-staged yet, so the audit starts from that.
- 2026-10-30: done-when rewritten as a spike (audit, format and lint comparison, decision); marked ready.
