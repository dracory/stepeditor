# Proposal: Predefined Condition Blocks

## Status
Proposed

## Description
Currently, all blocks in the Step Flow Editor must be defined by the user implementing the `CustomBlock` interface. While this provides maximum flexibility, common logic patterns like conditional branching (If/Else) and multi-way branching (Switch) are frequently needed across different implementations.

This proposal suggests introducing a set of "Predefined Blocks" that come built-in with the editor.

## Motivation
- **Standardization**: Provides a consistent way to handle logic across different workflows.
- **Ease of Use**: Users don't have to re-implement basic logic blocks for every project.
- **Enhanced UI**: Built-in blocks can have specialized UI components in the editor (e.g., a more intuitive condition builder).

## Proposed Changes

### 1. Built-in "If/Else" Block
A standard block with two branches: `True` and `False`.
- **Type**: `builtin:condition`
- **Default Data**: `variable`, `operator`, `value`.
- **Branches**: `["True", "False"]`.

### 2. Built-in "Switch" Block
A block that allows multiple named branches based on a variable's value.
- **Type**: `builtin:switch`
- **Default Data**: `variable`.
- **Dynamic Branches**: Ability for the user to add/remove branches in the editor UI.

### 3. Implementation Plan
- Modify the `Editor` struct to automatically include these predefined blocks in the `definitions` sent to the frontend.
- Update the Vue.js frontend to handle specialized rendering for these builtin types if necessary.
- Add a new `PredefinedBlocks` configuration option to enable/disable them.

## Alternatives Considered
- Keep them as examples: Users can copy-paste from documentation, but they lose out on specialized UI features.
- Provide a "Standard Library" package: Users can opt-in by adding `blockeditor.StandardConditions{}` to their `Blocks` config.
