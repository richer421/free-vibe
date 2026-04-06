# Console React Admin Shell Design

## Context

`templates/console-react` currently behaves like a lightweight Vite landing page template, while its local `AGENTS.md` defines it as a B-end admin frontend based on React, TypeScript, Tailwind CSS, and Ant Design.

The gap is not cosmetic. The template's current output does not express admin semantics, and it does not provide a usable starting point for backend-oriented business pages. This design closes that gap by turning `console-react` into a stable admin shell template without prematurely imposing routing, request, or auth architecture.

## Goal

Make `templates/console-react` a real admin-console starter that:

- keeps the stack as `React + TypeScript + Vite + Tailwind CSS + Ant Design`
- boots into a clear B-end interface instead of a marketing-style page
- provides a reusable layout shell and homepage structure
- remains intentionally light, so downstream projects can add routing, data, and auth without undoing the template

## Non-Goals

- no built-in router
- no request client abstraction
- no auth or permission model
- no global state framework
- no business-specific data schema

These are intentionally left open because the template should provide a strong shell, not lock projects into early architectural choices.

## Recommended Scope

The template should converge to a standard admin starting point:

1. install and wire `antd`
2. define a stable theme entry
3. replace the landing-page content with an admin shell
4. provide a default dashboard-style homepage
5. introduce only the minimum directory structure needed to keep responsibilities clear

This is the balance point between "too thin to be useful" and "too opinionated to extend."

## Visual Direction

Use a restrained, operational interface:

- light workspace background
- dark fixed sidebar for navigation
- compact top header for page title and actions
- dashboard-style content area with metrics, trend placeholder, quick actions, and recent activity

The template should look like an internal tool on first boot, not a product website.

The chosen direction is the "standard operations console" layout:

- fixed left navigation
- top header with page identity and action area
- central content region with reusable page spacing
- homepage blocks that demonstrate common admin information density

## Information Architecture

### Layout shell

The root interface should render three major zones:

- `Sider`: navigation and workspace identity
- `Header`: current page title, contextual description, and primary action slot
- `Content`: page body container with consistent spacing and section rhythm

The layout shell owns structure only. It should not know business data.

### Default homepage

The default homepage should demonstrate four reusable content patterns:

- overview stats
- trend or analytics placeholder
- quick actions
- recent activity

These blocks are intentionally generic. They teach structure without implying a specific business domain.

## Code Structure

Introduce a minimal but durable structure under `src/`:

- `layouts/`
  - `AdminLayout.tsx`: owns shell composition, nav items, header framing, and content container
- `pages/`
  - `Dashboard.tsx`: template homepage using reusable content blocks
- `components/`
  - shared presentational blocks such as stat cards, page sections, or quick action cards
- `theme/`
  - Ant Design theme tokens and theme helpers
- `index.css`
  - global foundation styles and Tailwind entry only

This split keeps layout, page content, reusable blocks, and theme concerns separate without creating a framework inside the template.

## Theming Strategy

Theme control should be centralized in Ant Design token configuration, with Tailwind used mainly for layout and light utility styling.

### Theme rules

- `antd` owns color tokens, radius, shadows, and component-level consistency
- Tailwind handles spacing, grid, flex layout, and selective visual accents
- global CSS remains thin and only covers document-level setup

### Default tone

- neutral-cool workspace surfaces
- dark navy/ink sidebar
- restrained blue accent for focus and primary actions
- moderate radius and soft shadowing for modern but stable admin visuals

This gives the template a professional baseline while remaining easy to rebrand.

## Dependency Expectations

The template should explicitly include the dependencies it claims to use.

Expected additions:

- `antd`
- `@ant-design/icons`

Expected app wiring:

- import Ant Design reset styles
- wrap the app with `ConfigProvider`
- provide theme tokens from a dedicated module

## Implementation Boundaries

### What the template should provide

- a running admin shell
- a real dashboard homepage
- reusable page container rhythm
- a small set of composable presentational primitives
- clear extension points for downstream features

### What the template should not provide

- router-specific code
- fetch helpers
- mock API contracts
- fake login flows
- speculative abstractions for future complexity

## Verification Criteria

The work is complete when all of the following are true:

- dependency installation succeeds
- `npm run lint` passes
- `npm run build` passes
- the default app renders as an admin dashboard shell rather than a landing page
- the resulting source structure is easy to extend without immediate refactor

## Risks And Mitigations

### Risk: visual over-design

If the template becomes too decorative, downstream teams will spend time removing aesthetics instead of starting work.

Mitigation:

- keep motion minimal
- prioritize hierarchy and readability over branding
- use admin-safe defaults

### Risk: premature architecture

If routing, auth, or requests are baked in now, the template will force assumptions that many projects will reject.

Mitigation:

- keep the template shell-only
- stop at structure and visual language

### Risk: mixed styling authority

If Tailwind and Ant Design both fully control component styling, maintenance becomes inconsistent.

Mitigation:

- make `antd` tokens the theme source of truth
- use Tailwind mainly for layout and page composition

## Rollout Notes

Implementation should prefer minimal, correct edits:

1. add missing dependencies
2. introduce theme config
3. replace `App.tsx` landing content with admin shell composition
4. add reusable presentational blocks only where they simplify the dashboard
5. verify lint and build before claiming completion

## Review Notes

The original brainstorming workflow calls for a subagent-based spec review loop. In this session, subagent delegation was not explicitly authorized by the user, so that step should be treated as a process deviation and replaced with local review unless the user opts in later.
