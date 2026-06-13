# Proposal: Workflow Validation

## Status
Proposed

## Description
Currently, the editor allows saving any flow, even if it contains incomplete configurations or invalid logic (e.g., a condition block with no branches).

This proposal suggests adding a validation layer that can be run both in the frontend (for immediate feedback) and in the backend (for security and integrity).

## Motivation
- **Reliability**: Ensure that only valid, runnable flows are saved to the database.
- **Developer Experience**: Provide clear error messages to the user about why their flow is invalid.

## Proposed Changes

### 1. Validation Rules
Define a set of standard rules:
- No empty branches.
- All required `Data` fields must be filled.
- No circular dependencies (if applicable).
- Custom validation rules provided by the `CustomBlock` implementation.

### 2. Go API Update
Add a `Validate` method to the `CustomBlock` or `BlockDefinition`.

```go
type CustomBlock interface {
    Definition() BlockDefinition
    Validate(data map[string]string) error
}
```

### 3. Frontend UI
- Display validation errors (e.g., a red icon or border) directly on the blocks that are failing.
- Disable the "Save" button or show a warning if the flow has validation errors.
- A "Validation Summary" panel listing all issues.

### 4. Backend Enforcement
- The `Editor.handleSave` method should run validation on the incoming JSON and return a `400 Bad Request` with error details if validation fails.

## Implementation Details
- Validation should be recursive, following the flow structure.
- We should provide a way to distinguish between "Errors" (must fix to save) and "Warnings" (can save, but might have issues).
