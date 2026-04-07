# Console React Admin Guardrails

## Overview

This skill keeps the `console-react` template aligned with a B-end admin-console baseline. It defines the default engineering constraints, directory boundaries, and UI-library usage rules for this template.

## Demo Interpretation

- The bootstrapped admin dashboard is a default demo, not a permanent product design baseline.
- Treat the demo as a reference for code style, module layout, and component composition speed only.
- Unless the user explicitly asks to preserve it, do not let the current demo's colors, copy, spacing rhythm, card layout, or dashboard information structure constrain later feature work.
- Reuse the scaffold shape, not the demo taste.

Current default demo files:

- `src/layouts/AdminLayout.tsx`
- `src/pages/dashboard/index.tsx`
- `src/pages/dashboard/components/ActivityTimeline.tsx`
- `src/pages/dashboard/components/QuickActionList.tsx`
- `src/config/app.ts`

## Core Rules

### React Hooks

- Do not use `useEffect` with an empty dependency array.
- Custom hooks must use the `useXXX` naming pattern.
- Each hook should own one focused responsibility.

### TypeScript

- Do not use `any`.
- Interface names must end with `Props`, `Model`, or `DTO`.
- Prefer explicit model types for module data instead of anonymous object arrays.

### Ant Design

- Tailwind is for layout and light decoration only; do not restyle Ant Design components by overriding their visual skin with utility classes.
- Prefer Ant Design's built-in controlled patterns and component state helpers.
- Use `message`, `modal`, and `notification` for global feedback.
- Do not create parallel popup abstractions when Ant Design already covers the case.

### Project Structure

Use module-first structure with a thin global layer:

```text
src/
├── pages/
│   └── <module>/
│       ├── index.tsx
│       ├── components/
│       ├── hooks/
│       └── service.ts
├── components/
├── hooks/
├── services/
├── utils/
├── config/
├── router/
└── types/
```

## Boundary Rules

- Module-private resources stay inside their own `pages/<module>/` directory.
- Cross-module logic moves to global directories only after reuse is real.
- Modules must not import each other's private components, hooks, or services.
- Keep shell structure, page modules, and reusable primitives clearly separated.

## Delivery Checklist

Before finishing work in this template, verify:

1. No `any`
2. No empty-deps `useEffect`
3. Custom hooks are named `useXXX`
4. Tailwind is not being used to repaint Ant Design components
5. New module code is contained under `pages/<module>/`
6. Shared logic is placed in the global layer only when genuinely shared
