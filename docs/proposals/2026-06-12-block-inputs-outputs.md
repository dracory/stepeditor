# Proposal: Block Inputs and Outputs

## Status
Proposed

## Description
Currently, blocks have a `Data` map for configuration, but there is no formal way to define data flowing *between* blocks. One block might produce an output (e.g., an HTTP response) that a subsequent block needs as an input.

This proposal suggests adding "Inputs" and "Outputs" to the `BlockDefinition` to facilitate data lineage and validation.

## Motivation
- **Type Safety**: Ensure that a block receiving data is getting the expected format.
- **Discoverability**: Help users see what data is available from previous steps.
- **Validation**: Prevent invalid connections or configurations where a required input is missing.

## Proposed Changes

### 1. Update `BlockDefinition`
Add `Inputs` and `Outputs` fields.

```go
type DataDefinition struct {
    Name string `json:"name"`
    Type string `json:"type"` // e.g., "string", "int", "json"
}

type BlockDefinition struct {
    // ... existing fields
    Inputs  []DataDefinition `json:"inputs,omitempty"`
    Outputs []DataDefinition `json:"outputs,omitempty"`
}
```

### 2. Update `Block` instance
Add a way for a block to "map" its inputs to outputs of previous blocks.

```go
type Block struct {
    // ... existing fields
    InputMappings map[string]string `json:"inputMappings,omitempty"` // Map local input name to "BlockID.OutputName"
}
```

### 3. Frontend Enhancements
- Visual indicators on blocks showing input/output ports.
- A "Data Picker" UI in the block settings to select outputs from preceding blocks in the flow.

## Implementation Details
- The editor needs to track the "scope" of available outputs at any given point in the recursive flow.
- Nested branches should have access to outputs from parent blocks.

## Impact
- **Breaking Change**: This might require updates to how `Block` and `BlockDefinition` are handled, though it can be done in a backward-compatible way by making the new fields optional.
