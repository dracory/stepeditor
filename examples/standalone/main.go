package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dracory/stepeditor"
)

// EmailStep defines a step for sending emails.
type EmailStep struct{}

func (b EmailStep) StepDefinition() stepeditor.StepDefinition {
	return stepeditor.StepDefinition{
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

// DelayStep defines a step for waiting.
type DelayStep struct{}

func (b DelayStep) StepDefinition() stepeditor.StepDefinition {
	return stepeditor.StepDefinition{
		Type:        "delay",
		Title:       "Wait",
		Description: "Wait for a specified duration.",
		Icon:        "bi-clock-fill",
		DefaultData: map[string]string{
			"duration": "10s",
		},
	}
}

// ConditionStep defines a step for conditional logic.
type ConditionStep struct{}

func (b ConditionStep) StepDefinition() stepeditor.StepDefinition {
	return stepeditor.StepDefinition{
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
	// Initialize the editor with custom steps
	editor := stepeditor.New(stepeditor.Config{
		Endpoint: "/editor",
		StepDefinitions: []stepeditor.CustomStep{
			EmailStep{},
			DelayStep{},
			ConditionStep{},
		},
		InitialFlow: []stepeditor.Step{
			{
				ID:    "s1",
				Type:  "delay",
				Title: "Daily Process",
				Data:  map[string]string{"duration": "24h"},
			},
			{
				ID:    "s2",
				Type:  "condition",
				Title: "Branch Check",
				Data: map[string]string{
					"variable": "status",
					"operator": "==",
					"value":    "approved",
				},
				Branches: map[string][]stepeditor.Step{
					"True": {
						{
							ID:    "s3",
							Type:  "email",
							Title: "Send Approval Email",
							Data:  map[string]string{"subject": "Approved!"},
						},
					},
					"False": {
						{
							ID:    "s4",
							Type:  "email",
							Title: "Send Rejection Email",
							Data:  map[string]string{"subject": "Rejected"},
						},
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
