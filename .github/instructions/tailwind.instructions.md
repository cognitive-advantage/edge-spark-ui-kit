---
description: 'Tailwind CSS v4 conventions and project theme'
applyTo: '**/*.html,**/*.js,**/*.css,**/tailwind.config.js'
---

# Tailwind CSS Instructions

## CRITICAL: This project uses Tailwind CSS v4

- Entry point: `layouts/default/input.css` uses `@import "tailwindcss"` (v4 syntax)
- Config: `@config "../../tailwind.config.js"` pulls in theme extensions
- CLI: `@tailwindcss/cli ^4.1.18` — NOT the legacy `tailwindcss` package
- Build: `npm run build:css` or `npm run watch:css`
- Output: `layouts/default/css/output.css`

## Project Color Palette

The ONLY colors defined in `tailwind.config.js` are:

| Token    | Purpose                        |
| -------- | ------------------------------ |
| primary  | Brand / interactive elements   |
| red      | Errors, destructive actions    |
| salmon   | Warnings, in-progress states   |
| lime     | Success, complete states       |
| grey     | Neutral text, borders, backgrounds |

## Rules

1. **ONLY use project theme colors** — `primary`, `red`, `salmon`, `lime`, `grey`.
2. **NEVER use default Tailwind colors** that are not in the theme — no `green-*`, `blue-*`, `amber-*`, `yellow-*`, `indigo-*`, `teal-*`, etc.
3. **Color mapping**:
   - Success/complete → `lime` (NOT `green`)
   - Active/interactive → `primary` (NOT `blue`)
   - Warning/in-progress → `salmon` (NOT `amber`/`yellow`)
   - Error/destructive → `red`
   - Neutral → `grey`
4. **Rebuild CSS** after adding new utility classes: `npm run build:css` or ensure `watch:css` is running.
5. **Content scanning** is configured to scan `layouts/**/*.html` and `layouts/**/*.js`.
