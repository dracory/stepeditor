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
		Title:       "Branch Check",
		Description: "Check a condition and branch the flow.",
		Icon:        "bi-shuffle",
		DefaultData: map[string]string{
			"variable": "status",
			"operator": "==",
			"value":    "approved",
		},
		BranchNames: []string{"True", "False"},
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

	// Add an initial branched flow
	editor.SetFlow([]stepfloweditor.Block{
		{
			ID:    "b1",
			Type:  "delay",
			Title: "Daily Process",
			Data:  map[string]string{"duration": "24h"},
		},
		{
			ID:    "b2",
			Type:  "condition",
			Title: "Branch Check",
			Data: map[string]string{
				"variable": "status",
				"operator": "==",
				"value":    "approved",
			},
			Branches: map[string][]stepfloweditor.Block{
				"True": {
					{
						ID:    "b3",
						Type:  "email",
						Title: "Send Approval Email",
						Data:  map[string]string{"subject": "Approved!"},
					},
				},
				"False": {
					{
						ID:    "b4",
						Type:  "email",
						Title: "Send Rejection Email",
						Data:  map[string]string{"subject": "Rejected"},
					},
				},
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
