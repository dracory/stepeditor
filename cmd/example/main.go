package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/username/stepfloweditor"
)

// EmailBlock defines a block for sending emails.
type EmailBlock struct{}

func (b EmailBlock) Definition() stepfloweditor.BlockDefinition {
	return stepfloweditor.BlockDefinition{
		Type:        "email",
		Title:       "Send Email",
		Description: "Sends an email to a recipient.",
		Icon:        "bi-envelope-fill",
		DefaultData: map[string]string{
			"to":      "",
			"subject": "",
			"body":    "",
		},
	}
}

// DelayBlock defines a block for waiting.
type DelayBlock struct{}

func (b DelayBlock) Definition() stepfloweditor.BlockDefinition {
	return stepfloweditor.BlockDefinition{
		Type:        "delay",
		Title:       "Wait",
		Description: "Wait for a specified duration.",
		Icon:        "bi-clock-fill",
		DefaultData: map[string]string{
			"duration": "10s",
		},
	}
}

// ConditionBlock defines a block for conditional logic.
type ConditionBlock struct{}

func (b ConditionBlock) Definition() stepfloweditor.BlockDefinition {
	return stepfloweditor.BlockDefinition{
		Type:        "condition",
		Title:       "Condition",
		Description: "Check a condition before proceeding.",
		Icon:        "bi-question-diamond-fill",
		DefaultData: map[string]string{
			"expression": "",
		},
	}
}

func main() {
	// Initialize the editor with custom blocks
	editor := stepfloweditor.New(stepfloweditor.NewConfig{
		Endpoint: "/editor",
		Blocks: []stepfloweditor.CustomBlock{
			EmailBlock{},
			DelayBlock{},
			ConditionBlock{},
		},
	})

	// Add some initial blocks to the flow
	editor.SetFlow([]stepfloweditor.Block{
		{
			ID:    "init_1",
			Type:  "email",
			Title: "Welcome Email",
			Data: map[string]string{
				"to":      "user@example.com",
				"subject": "Welcome!",
				"body":    "Hello and welcome to our service.",
			},
		},
		{
			ID:    "init_2",
			Type:  "delay",
			Title: "Wait for 1 day",
			Data: map[string]string{
				"duration": "24h",
			},
		},
	})

	// Mount the editor on the specified endpoint
	http.Handle("/editor/", editor)

	// Add a simple landing page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, `<html><body>
			<h1>Step Flow Editor Example</h1>
			<p>Go to <a href="/editor/">the editor</a> to build your flow.</p>
		</body></html>`)
	})

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Editor available at http://localhost:8080/editor/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
