# Console React Admin Shell Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn `templates/console-react` from a landing-page starter into a usable Ant Design admin shell template with a dashboard homepage and minimal test coverage.

**Architecture:** Keep the template shell-only. Use Ant Design tokens as the theme source of truth, Tailwind for layout and composition, and split the code into `theme`, `layouts`, `pages`, and small presentational `components`. Add a minimal Vitest + Testing Library setup so the template has executable UI smoke tests without introducing a broader app architecture.

**Tech Stack:** React 19, TypeScript, Vite, Tailwind CSS 4, Ant Design 5, Vitest, Testing Library

---

## Chunk 1: Dependencies And Test Harness

### Task 1: Add admin-shell and test dependencies

**Files:**
- Modify: `templates/console-react/package.json`

- [ ] **Step 1: Write the failing test dependency expectations**

Add a test script target and the dependency list required for:
- `antd`
- `@ant-design/icons`
- `vitest`
- `jsdom`
- `@testing-library/react`
- `@testing-library/jest-dom`

The first red state is "test command cannot run because the script and packages are missing."

- [ ] **Step 2: Run the missing test command to verify the red state**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test`

Expected: command fails because `test` script does not exist.

- [ ] **Step 3: Add the minimal dependency and script wiring**

Update `templates/console-react/package.json` to:
- add `test` script using `vitest run`
- add `test:watch` script using `vitest`
- add the Ant Design runtime dependencies
- add the Vitest and Testing Library dev dependencies

- [ ] **Step 4: Run the test command again to verify the next red state**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm install && npm test`

Expected: Vitest starts, then fails because no test files or config exist yet.

- [ ] **Step 5: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/package.json
git commit -m "test: add console-react test dependencies"
```

### Task 2: Create the template test harness

**Files:**
- Modify: `templates/console-react/vite.config.ts`
- Create: `templates/console-react/src/test/setup.ts`

- [ ] **Step 1: Write the failing render test harness expectation**

Plan for a browser-like test environment with:
- `jsdom`
- global setup loading `@testing-library/jest-dom`

The red state is "the future component test cannot run in a DOM environment."

- [ ] **Step 2: Run the existing test command to confirm no harness exists yet**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test`

Expected: FAIL due to missing test files or missing Vitest setup.

- [ ] **Step 3: Add the minimal Vitest configuration**

Update `templates/console-react/vite.config.ts` to add:

```ts
test: {
  environment: 'jsdom',
  setupFiles: './src/test/setup.ts',
}
```

Create `templates/console-react/src/test/setup.ts` with:

```ts
import '@testing-library/jest-dom'
```

- [ ] **Step 4: Run tests to confirm the harness is ready for component tests**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test`

Expected: Vitest runs successfully but reports no tests found.

- [ ] **Step 5: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/vite.config.ts /Users/richer/richer/free-vibe-coding/templates/console-react/src/test/setup.ts
git commit -m "test: configure console-react vitest harness"
```

## Chunk 2: Theme Wiring And Shell Composition

### Task 3: Add a failing app-level smoke test

**Files:**
- Create: `templates/console-react/src/App.test.tsx`

- [ ] **Step 1: Write the failing test**

Create `templates/console-react/src/App.test.tsx`:

```tsx
import { render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders the admin shell dashboard', () => {
    render(<App />)

    expect(screen.getByRole('heading', { name: 'Operational Overview' })).toBeInTheDocument()
    expect(screen.getByText('Workspace')).toBeInTheDocument()
    expect(screen.getByText('Quick Actions')).toBeInTheDocument()
    expect(screen.getByText('Recent Activity')).toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run the test to verify it fails for the right reason**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test -- --runInBand`

Expected: FAIL because the current landing page does not render the admin shell content.

- [ ] **Step 3: Keep the failing test unchanged and move to implementation**

Do not weaken the assertions. They capture the required shell semantics.

- [ ] **Step 4: Commit the failing test**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/src/App.test.tsx
git commit -m "test: define console-react admin shell smoke test"
```

### Task 4: Add theme tokens and app providers

**Files:**
- Create: `templates/console-react/src/theme/themeConfig.ts`
- Modify: `templates/console-react/src/main.tsx`
- Modify: `templates/console-react/src/index.css`

- [ ] **Step 1: Implement the minimal theme module**

Create `templates/console-react/src/theme/themeConfig.ts` exporting an Ant Design theme config with:
- cool neutral surfaces
- dark sider colors
- restrained blue primary color
- stable radius and shadow values

Use a plain exported object, not a custom theme abstraction.

- [ ] **Step 2: Wire the provider**

Update `templates/console-react/src/main.tsx` to wrap `<App />` in:

```tsx
<ConfigProvider theme={themeConfig}>
  <App />
</ConfigProvider>
```

- [ ] **Step 3: Tighten the global stylesheet**

Keep `templates/console-react/src/index.css` focused on:
- Tailwind import
- base font stack
- root/document sizing
- background and text defaults aligned with the dashboard shell

Do not place page-specific styles here.

- [ ] **Step 4: Run the smoke test**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test -- --runInBand`

Expected: still FAIL, but now only because the admin shell components are not implemented yet.

- [ ] **Step 5: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/src/theme/themeConfig.ts /Users/richer/richer/free-vibe-coding/templates/console-react/src/main.tsx /Users/richer/richer/free-vibe-coding/templates/console-react/src/index.css
git commit -m "feat: add console-react theme provider"
```

### Task 5: Implement the admin shell layout

**Files:**
- Create: `templates/console-react/src/layouts/AdminLayout.tsx`
- Modify: `templates/console-react/src/App.tsx`

- [ ] **Step 1: Implement `AdminLayout.tsx`**

Create a focused layout component that renders:
- fixed-width dark `Sider`
- workspace identity and simple nav items
- top `Header` with eyebrow label, title, supporting copy, and a primary action button
- `Content` slot rendered via `children`

Use Ant Design `Layout`, `Menu`, `Avatar`, `Button`, and small iconography. Use Tailwind classes only for layout composition and spacing.

- [ ] **Step 2: Make `App.tsx` compose the shell**

Update `templates/console-react/src/App.tsx` to render:

```tsx
import { AdminLayout } from './layouts/AdminLayout'
import { Dashboard } from './pages/Dashboard'

function App() {
  return (
    <AdminLayout>
      <Dashboard />
    </AdminLayout>
  )
}
```

- [ ] **Step 3: Run the smoke test**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test -- --runInBand`

Expected: FAIL, but only on dashboard content that has not been added yet.

- [ ] **Step 4: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/src/layouts/AdminLayout.tsx /Users/richer/richer/free-vibe-coding/templates/console-react/src/App.tsx
git commit -m "feat: add console-react admin layout shell"
```

## Chunk 3: Dashboard Page And Reusable Blocks

### Task 6: Implement reusable dashboard presentation blocks

**Files:**
- Create: `templates/console-react/src/components/StatCard.tsx`
- Create: `templates/console-react/src/components/PageSection.tsx`

- [ ] **Step 1: Implement `StatCard.tsx`**

Create a small presentational card component with props for:
- label
- value
- helper text
- optional accent style

Keep it display-only.

- [ ] **Step 2: Implement `PageSection.tsx`**

Create a simple wrapper that standardizes section title, extra slot, and content padding.

- [ ] **Step 3: Run tests**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test -- --runInBand`

Expected: still FAIL because the `Dashboard` page is not wired yet.

- [ ] **Step 4: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/src/components/StatCard.tsx /Users/richer/richer/free-vibe-coding/templates/console-react/src/components/PageSection.tsx
git commit -m "feat: add console-react dashboard building blocks"
```

### Task 7: Implement the dashboard homepage

**Files:**
- Create: `templates/console-react/src/pages/Dashboard.tsx`

- [ ] **Step 1: Implement the page**

Create `templates/console-react/src/pages/Dashboard.tsx` that renders:
- a four-card stats row
- a trend placeholder section
- a quick actions section
- a recent activity section

Make sure the page includes these visible strings for the smoke test:
- `Operational Overview`
- `Workspace`
- `Quick Actions`
- `Recent Activity`

Prefer static seed data arrays in the file over premature extraction.

- [ ] **Step 2: Run the smoke test to verify green**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test -- --runInBand`

Expected: PASS

- [ ] **Step 3: Refactor only if needed**

If the page feels too large after implementation, extract one more tiny presentational helper. Do not add abstractions preemptively.

- [ ] **Step 4: Commit**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react/src/pages/Dashboard.tsx
git commit -m "feat: add console-react dashboard homepage"
```

## Chunk 4: Verification And Template Fit

### Task 8: Run full verification

**Files:**
- Verify only: `templates/console-react/package.json`
- Verify only: `templates/console-react/src/App.tsx`
- Verify only: `templates/console-react/src/layouts/AdminLayout.tsx`
- Verify only: `templates/console-react/src/pages/Dashboard.tsx`
- Verify only: `templates/console-react/src/theme/themeConfig.ts`

- [ ] **Step 1: Run lint**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm run lint`

Expected: PASS

- [ ] **Step 2: Run tests**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm test`

Expected: PASS

- [ ] **Step 3: Run production build**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm run build`

Expected: PASS

- [ ] **Step 4: Manual template review**

Run: `cd /Users/richer/richer/free-vibe-coding/templates/console-react && npm run dev`

Verify manually:
- first screen is clearly an admin dashboard
- sidebar, header, and content spacing feel coherent
- homepage cards and sections read as reusable business UI, not demo marketing content

- [ ] **Step 5: Commit the verified implementation**

```bash
git add /Users/richer/richer/free-vibe-coding/templates/console-react
git commit -m "feat: converge console-react into admin shell template"
```

## Notes For Execution

- Keep edits local to `templates/console-react` unless verification exposes a repository-level issue.
- Do not change routing, auth, request, or state-management architecture.
- Prefer one focused component per file.
- If Ant Design class names and Tailwind utilities conflict, simplify the Tailwind usage instead of fighting the component library.
- `.superpowers/` is local brainstorming output and should not be staged as part of implementation.

## Review Notes

The original plan-writing workflow expects a plan-reviewer subagent loop. In this session, explicit user authorization for subagent delegation was not provided, so local review should be used instead of agent delegation.
