package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/dracory/stepeditor"
)

type MyBlock struct{}

func (b MyBlock) Definition() blockeditor.BlockDefinition {
	return blockeditor.BlockDefinition{
		Type:        "myblock",
		Title:       "My Custom Block",
		Description: "A block for the embedded editor.",
		Icon:        "bi-star-fill",
	}
}

func main() {
	// 1. Initialize the editor
	editor := blockeditor.New(blockeditor.NewConfig{
		ID:       "my-embedded-editor",
		Endpoint: "/api/editor",
		Blocks:   []blockeditor.CustomBlock{MyBlock{}},
	})

	// 2. Setup the API endpoint for the editor
	http.Handle("/api/editor/", editor)

	// 3. Render the editor in a custom layout
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Get the editor's HTML
		editorHTML := template.HTML(editor.ToHTML())

		layout := `
<!DOCTYPE html>
<html>
<head>
    <title>Embedded Editor Example</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body { padding: 2rem; background: #eee; }
        .editor-container {
            background: white;
            padding: 1rem;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>My Custom App</h1>
        <p>This page demonstrates embedding the Step Flow Editor as a component.</p>

        <div class="editor-container">
            {{.Editor}}
        </div>

        <div class="mt-4">
            <button class="btn btn-secondary">External Action</button>
        </div>
    </div>
</body>
</html>`

		tmpl := template.Must(template.New("layout").Parse(layout))
		tmpl.Execute(w, map[string]interface{}{
			"Editor": editorHTML,
		})
	})

	fmt.Println("Server starting on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
