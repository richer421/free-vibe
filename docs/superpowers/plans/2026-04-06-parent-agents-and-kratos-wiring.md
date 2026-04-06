# Parent AGENTS And Kratos Wiring Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `freevibe` bootstrap a parent-repo-level `AGENTS.md` and close the `kratos` template's local-skill wiring gap.

**Architecture:** Keep the change small and explicit. Generate the parent repo `AGENTS.md` from scaffold code alongside the root `README.md` and `Makefile`, and update the checked-in `templates/kratos/AGENTS.md` so it always routes work through the existing local proxy skill before dropping into downstream Kratos skills.

**Tech Stack:** Go, Go tests, Markdown

---

## Chunk 1: Parent Repo AGENTS Bootstrap

### Task 1: Add tests for parent repo AI guidance generation

**Files:**
- Modify: `internal/freevibe/scaffold/scaffold_test.go`

- [ ] **Step 1: Write a failing test**

Add a test that initializes a parent project and asserts the generated root `AGENTS.md` exists and contains submodule-orchestration guidance.

- [ ] **Step 2: Run the focused test to verify red**

Run: `go test ./internal/freevibe/scaffold -run TestInitProjectCreatesRootAgentsFile`

Expected: FAIL because the file is not generated yet.

- [ ] **Step 3: Implement the minimal scaffold change**

Add a helper that writes a default root `AGENTS.md` during parent project initialization.

- [ ] **Step 4: Re-run the focused test**

Run: `go test ./internal/freevibe/scaffold -run TestInitProjectCreatesRootAgentsFile`

Expected: PASS.

## Chunk 2: Kratos Template Skill Wiring

### Task 2: Add tests for the checked-in Kratos template guidance

**Files:**
- Modify: `internal/freevibe/scaffold/scaffold_test.go`
- Modify: `templates/kratos/AGENTS.md`

- [ ] **Step 1: Write a failing test**

Add a test that reads `templates/kratos/AGENTS.md` from the repository and asserts it references the existing local proxy skill and the downstream Kratos local skills.

- [ ] **Step 2: Run the focused test to verify red**

Run: `go test ./internal/freevibe/scaffold -run TestKratosTemplateAgentsConnectsLocalSkills`

Expected: FAIL because the current file does not mention the local skills.

- [ ] **Step 3: Implement the minimal documentation change**

Update `templates/kratos/AGENTS.md` so work inside the template always starts with the local proxy skill and continues into the existing Kratos local skills when routed there.

- [ ] **Step 4: Run the package test suite**

Run: `go test ./internal/freevibe/scaffold`

Expected: PASS.
