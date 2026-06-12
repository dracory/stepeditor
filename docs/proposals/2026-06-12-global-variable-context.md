# Proposal: Global Variable and Context Management

## Status
Proposed

## Description
In complex workflows, certain data needs to be accessible everywhere, regardless of block nesting or flow position. This includes things like user IDs, environment flags, or session tokens.

This proposal suggests a way to define and manage "Global Variables" or "Flow Context" within the editor.

## Motivation
- **Simplicity**: Avoid passing the same data through every single block as an input/output.
- **Configurability**: Allow users to define constants or environment-specific values at the flow level.

## Proposed Changes

### 1. New Config Option
Add a `Context` or `Variables` field to `NewConfig`.

```go
type NewConfig struct {
    // ...
    Variables []VariableDefinition `json:"variables"`
}
```

### 2. Editor UI
- Add a "Global Variables" tab or modal to the editor.
- Allow users to define name, type, and default value for these variables.

### 3. Usage in Blocks
- Update the expression language (if any) or block data fields to allow referencing these variables (e.g., `{{global.user_id}}`).

## Implementation Plan
- Store global variables in the top-level flow JSON.
- Provide a helper function in Go to inject initial variables.
- Update frontend to provide a centralized place to manage these.

## Alternatives
- Just use a special "Set Variable" block. While useful, it doesn't provide a way to see all available variables in one place or define initial values.
