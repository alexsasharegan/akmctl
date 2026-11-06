# Agent Environment Management

## Context

Agent environments consist of more than installed skills. They include repository instructions, reusable workflows, documentation, tools, conventions, guardrails, and the mechanisms that connect them.

These components evolve independently but must work together coherently. This creates an environment-management problem that conventional package management does not fully address.

A related challenge is organizing knowledge so humans and agents can discover, understand, adopt, operate, troubleshoot, and extend capabilities.

Agent environment management should therefore address both the composition of the environment and the lifecycle of the knowledge that makes it useful.

## Design Principles

### 1. Separate Ownership

Distinguish upstream repository instructions from locally managed additions. Integrations should be independently identifiable, updateable, and removable without modifying or overwriting upstream content.

### 2. Compose Instructions

Treat instructions as layered configuration rather than a monolithic `AGENTS.md`.

Potential layers include:

- Universal agent instructions.
- Repository-owned instructions.
- Locally managed repository adaptations.
- Skill-specific instructions.

Keep entry-point instructions concise and delegate specialized guidance to discoverable documentation and skills.

### 3. Manage Dependencies

Skills depend on more than other skills. They may require particular files, conventions, tools, agent capabilities, or repository structures.

Make these dependencies explicit enough to identify incompatibilities and coordinate changes.

### 4. Treat Migrations as First-Class

Skill and instruction updates must account for additions, modifications, renames, removals, and breaking changes.

A successful update reconciles the installed environment with the declared desired state. Migration guidance should explain compatibility requirements and provide actionable steps.

### 5. Optimize Instruction Placement

Place information where and when it is needed.

- Mechanical rules belong in automated checks.
- Judgment-based guidance belongs in instructions.
- Repeatable procedures belong in skills.
- Architectural context belongs in documentation.
- Navigation guidance belongs in repository orientation materials.

Instructions should be concise, discoverable, and nonredundant. More instructions do not necessarily produce better agent behavior.

### 6. Treat Capability Guides as First-Class Knowledge Artifacts

A capability guide explains how to understand, adopt, operate, troubleshoot, and extend a capability.

It serves both humans and agents, connecting conceptual explanations with actionable knowledge.

A useful guide follows the capability adoption journey:

- **Overview:** What the capability is, what it does, and why it exists.
- **How It Works:** Conceptual model, architecture, and terminology.
- **Getting Started:** Prerequisites, installation, and configuration.
- **Walkthrough:** A practical end-to-end example.
- **Validation:** How to test and verify expected behavior.
- **Troubleshooting:** Common problems, diagnostics, and recovery.
- **Distribution:** How to share, update, and manage the capability.
- **Examples and Applications:** Additional use cases and extension ideas.
- **Further Reading:** Links to authoritative references and related material.

Not every guide requires every section. Its structure should follow the adoption journey rather than impose a documentation checklist.

Capability guides complement, rather than replace:

- Skills, which provide operational procedures.
- Instructions, which establish behavioral expectations.
- Reference documentation, which defines authoritative interfaces and schemas.
- Changelogs, which communicate changes between versions.

Knowledge should be organized around successful capability adoption, not merely around the structure of the software being documented.

### 7. Continuously Improve the Environment

Treat agent sessions as opportunities to improve the environment itself.

Retrospectives should identify recurring friction, missing context, ineffective instructions, automation opportunities, and unnecessary complexity. Improvements should target the underlying cause rather than reflexively adding more instructions.

## Architectural Questions

1. How should independently owned instruction sets be composed within an existing repository?
2. How should skills declare and validate their environmental dependencies?
3. How should updates, removals, and breaking changes propagate across skills and repository conventions?
4. How can instruction composition remain compatible across different agent harnesses?
5. How should capability knowledge be structured, maintained, versioned, and connected to operational instructions, skills, and authoritative reference documentation?
6. How can observations from agent sessions become actionable improvements to the environment?
7. How can instruction and documentation effectiveness be evaluated rather than assumed?

## Guiding Principle

**Manage the agent environment as a versioned, composable, and continuously evolving system—not merely as a collection of prompts and installed skills.**

The objective is to make agent behavior more reliable by improving the environment in which agents operate, while preserving clear ownership boundaries, making knowledge discoverable, and minimizing instructional overhead.
