# Proposal: Naming Standardization for Step Flow Editor

## Motivation
The current codebase uses a mix of "Block" and "Step" terminology. While the project is named "Step Flow Editor", many internal types and the Go package are named "blockeditor". This proposal aims to standardize all naming around the concept of "Steps" and "Flows" to improve consistency and clarity.

## Proposed Changes

### 1. Go Package and Module
- **Module Path**: Remains `github.com/dracory/stepeditor` (already matches).
- **Package Name**: Rename from `blockeditor` to `stepeditor`.

### 2. Core Types (Go)
| Current Name | Proposed Name |
|--------------|---------------|
| `Block` | `Step` |
| `BlockDefinition` | `StepDefinition` |
| `CustomBlock` | `CustomStep` |
| `NewConfig` | `Config` |

### 3. Struct Fields (Go)
**`stepeditor.Config` (formerly `NewConfig`)**:
| Current Field | Proposed Field |
|---------------|----------------|
| `Value []Block` | `InitialFlow []Step` |
| `Blocks []CustomBlock` | `StepDefinitions []CustomStep` |

**`stepeditor.Step` (formerly `Block`)**:
- No major field changes besides types.

**`stepeditor.StepDefinition` (formerly `BlockDefinition`)**:
- `BranchNames []string` remains the same.

### 4. Methods (Go)
| Current Method | Proposed Method |
|----------------|-----------------|
| `CustomStep.Definition()` | `CustomStep.StepDefinition()` |
| `Editor.GetFlow()` | `Editor.GetFlow()` (returns `[]Step`) |
| `Editor.SetFlow()` | `Editor.SetFlow()` (takes `[]Step`) |
| `Editor.ToHTML()` | `Editor.ToHTML()` (remains same) |

### 5. Frontend (JavaScript/Vue)
- **Component**: `flow-node` -> `step-node`
- **JS Variable Names**:
    - `block` -> `step`
    - `definitions` -> `stepDefinitions`
    - `addBlock` -> `addStep`
    - `addBlockToBranch` -> `addStepToBranch`
    - `selectBlock` -> `selectStep`
    - `removeBlock` -> `removeStep`

### 6. Frontend (CSS)
Rename CSS classes to use CamelCase and "Step" instead of "Block" to avoid conflicts with global libraries like Bootstrap.

| Current Class | Proposed Class |
|---------------|----------------|
| `.editor-root` | `.EditorRoot` |
| `.editor-header` | `.EditorHeader` |
| `.editor-body` | `.EditorBody` |
| `.toolbox` | `.Toolbox` |
| `.toolbox-header` | `.ToolboxHeader` |
| `.block-library-item` | `.StepLibraryItem` |
| `.canvas` | `.Canvas` |
| `.settings-panel` | `.SettingsPanel` |
| `.flow-block-container` | `.FlowStepContainer` |
| `.block-card` | `.StepCard` |
| `.block-icon` | `.StepIcon` |
| `.branches-container` | `.BranchesContainer` |
| `.branch` | `.Branch` |
| `.branch-label` | `.BranchLabel` |
| `.connector-v` | `.ConnectorV` |
| `.block-actions` | `.StepActions` |

### 7. Documentation and Examples
- Update `README.md` to use the new package name and types.
- Update all example files in `examples/` to reflect these changes.

## Impact
This is a breaking change for users of the library as it renames the primary package and core configuration structures. However, it will result in a much more intuitive API that matches the project's stated purpose.
