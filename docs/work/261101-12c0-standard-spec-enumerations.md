---
id: "261101-12c0"
title: standard spec enumerations
summary: Audit existing conventions and external prior art to establish a consistent approach to enumerating specification rules.
status: raw
kind: spike
created: 2026-11-01
updated: 2026-11-01
---

# specification rule enumerations

## Goal

Establish a principled convention for enumerating and classifying specification items across the repository.

The `check` command currently specifies rules using identifiers such as WK, I, A, RM, G, L, S, CL, and SF. The rationale and consistency of these designations are unclear.

Conduct a prior-art audit to understand how established engineering specifications identify, organize, classify, and express rules. Use the findings to recommend a consistent convention for this project.

## Investigation

### Repository audit

Survey existing specification enumerations throughout the repository, particularly those used by the `check` command.

Identify what each enumeration represents, how it is used, and whether different conventions overlap or conflict.

### External prior art

Survey representative specifications from both established standards bodies and modern engineering projects.

Potential examples:

- **IETF:** HTTP (RFC 9110), TLS, JSON (RFC 8259), RFC 2119.
- **IEEE / Ecma:** POSIX, ECMAScript.
- **W3C / WHATWG:** WCAG, CSS, HTML.
- **Modern projects:** Rust RFCs, Kubernetes KEPs, OpenTelemetry, OCI, JSON Schema.

Investigate how they handle:

- Rule identifiers and hierarchical numbering.
- Requirement classifications and normative keywords (MUST, SHOULD, MAY).
- Conformance requirements and compliance criteria.
- Cross-references, exceptions, and evolution.

### Comparative analysis

Compare the approaches and identify useful conventions, trade-offs, and recurring patterns.

Pay particular attention to the distinction between:

- **Identity:** How an individual rule is referenced.
- **Classification:** What kind of rule it is.
- **Normativity:** How strongly it is required.

## Deliverables

- Inventory of existing repository conventions.
- Curated prior-art references.
- Comparative analysis of relevant approaches.
- Recommended convention with illustrative examples.
- Migration considerations for existing enumerations.

## Done When

- [ ] Existing repository enumerations have been surveyed.
- [ ] Relevant standards-body and modern engineering examples have been examined.
- [ ] Identity, classification, and normative force have been evaluated separately.
- [ ] A recommended convention has been proposed with concrete examples.
- [ ] Trade-offs and migration considerations are documented.

## Log

- 2026-11-01: filed.
