# Proposal: Undo/Redo Support

## Status
Proposed

## Description
Editing complex flows can be error-prone. Accidentally deleting a branch or moving a block incorrectly can be frustrating if there is no way to quickly revert.

This proposal suggests implementing an Undo/Redo mechanism in the frontend editor.

## Motivation
- **User Confidence**: Users are more likely to experiment and build complex flows if they know they can easily revert changes.
- **Productivity**: Fixing a mistake is much faster with a single keyboard shortcut than manually re-configuring a block.

## Proposed Changes

### 1. Frontend State Management
- Implement a history stack in the Vue.js application.
- After every significant action (add block, delete block, move block, update data), push a snapshot of the current flow to the stack.

### 2. UI Elements
- Add "Undo" and "Redo" buttons to the editor toolbar.
- Support standard keyboard shortcuts: `Ctrl+Z` (Undo) and `Ctrl+Y` / `Ctrl+Shift+Z` (Redo).

### 3. Optimization
- Instead of full snapshots, consider using JSON patches to save memory and improve performance for very large flows.

## Implementation Details
- The history should be local to the session and doesn't necessarily need to be persisted to the backend unless "Auto-save" is also implemented.
- Limit the history depth (e.g., 50 actions) to prevent excessive memory usage.

## Impact
- purely frontend change, no changes required to the Go backend structures.
