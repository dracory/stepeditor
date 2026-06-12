# Step Flow Editor

A step-by-step workflow editor for Golang projects, similar to Power Automate. This library provides an embeddable UI that allows users to build and configure sequential flows using custom-defined blocks.

![Step Flow Editor Screenshot](assets/screenshot.png)

## Features

- **Embeddable**: Easily mount the editor on any HTTP endpoint in your Go application.
- **Extensible**: Define your own block types by implementing a simple Go interface.
- **Reactive UI**: Built with Vue.js 3 and Bootstrap for a modern, responsive experience.
- **JSON Serialization**: Load and save flows as JSON arrays.
- **Interactive**: Add, remove, reorder, and configure blocks with ease.

## Installation

```bash
go get github.com/username/stepfloweditor
```

## Quick Start

```go
package main

import (
	"net/http"
	"github.com/username/stepfloweditor"
)

// Define a custom block
type MyBlock struct{}

func (b MyBlock) Definition() stepfloweditor.BlockDefinition {
	return stepfloweditor.BlockDefinition{
		Type:        "my-block",
		Title:       "My Custom Block",
		Description: "Does something awesome.",
		Icon:        "bi-star-fill",
		DefaultData: map[string]string{
			"setting": "default value",
		},
	}
}

func main() {
	editor := stepfloweditor.New(stepfloweditor.NewConfig{
		Endpoint: "/editor",
		Blocks: []stepfloweditor.CustomBlock{
			MyBlock{},
		},
	})

	// Mount the editor
	http.Handle("/editor/", editor)

	http.ListenAndServe(":8080", nil)
}
```

## API

### `stepfloweditor.New(config NewConfig) *Editor`

Creates a new editor instance.

### `Editor.ServeHTTP(w, r)`

Handles HTTP requests. Mount this on your router. Ensure the path ends with a trailing slash or is handled correctly by your router.

### `Editor.GetFlow() []Block`

Returns the current flow as a slice of `Block` structs.

### `Editor.SetFlow(flow []Block)`

Sets the current flow.

## Custom Blocks

To create a custom block, implement the `CustomBlock` interface:

```go
type CustomBlock interface {
	Definition() BlockDefinition
}
```

The `BlockDefinition` includes:
- `Type`: Unique identifier for the block type.
- `Title`: Display name in the library and editor.
- `Description`: Short description of what the block does.
- `Icon`: Bootstrap Icon class (e.g., `bi-envelope-fill`).
- `DefaultData`: Map of default attributes for the block.
