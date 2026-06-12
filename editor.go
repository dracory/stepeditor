package blockeditor

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	texttemplate "text/template"
)

//go:embed internal/assets/*
var assets embed.FS

var templates = texttemplate.Must(texttemplate.New("").Delims("[[", "]]").ParseFS(assets, "internal/assets/*.html"))

// Block represents a single step in the flow.
type Block struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Data     map[string]string `json:"data"`
	Branches map[string][]Block `json:"branches,omitempty"`
}

// BlockDefinition provides information about a block type.
type BlockDefinition struct {
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"` // Bootstrap icon class
	DefaultData map[string]string `json:"defaultData"`
	BranchNames []string          `json:"branchNames,omitempty"`
}

// CustomBlock is an interface that users can implement to provide their own blocks.
type CustomBlock interface {
	Definition() BlockDefinition
}

// NewConfig is the configuration for the Step Flow Editor.
type NewConfig struct {
	ID       string
	Endpoint string
	Value    []Block
	Blocks   []CustomBlock
}

// Editor handles the HTTP requests for the step flow editor.
type Editor struct {
	id     string
	config NewConfig
	flow   []Block
	mu     sync.RWMutex
}

// New creates a new Step Flow Editor.
func New(config NewConfig) *Editor {
	// Ensure endpoint doesn't end with slash for consistency
	config.Endpoint = strings.TrimSuffix(config.Endpoint, "/")

	id := config.ID
	if id == "" {
		b := make([]byte, 4)
		rand.Read(b)
		id = fmt.Sprintf("editor_%x", b)
	}

	if config.Value == nil {
		config.Value = make([]Block, 0)
	}

	return &Editor{
		id:     id,
		config: config,
		flow:   config.Value,
	}
}

// ServeHTTP implements the http.Handler interface.
func (e *Editor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, e.config.Endpoint)

	switch {
	case path == "" || path == "/":
		e.serveIndex(w, r)
	case path == "/config":
		e.serveConfig(w, r)
	case path == "/save" && r.Method == http.MethodPost:
		e.handleSave(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (e *Editor) serveIndex(w http.ResponseWriter, r *http.Request) {
	data := e.getTemplateData()

	// Use text/template for JS and CSS to avoid auto-escaping
	tJS := texttemplate.Must(texttemplate.New("js").Delims("[[", "]]").Parse(string(data.JS)))
	var jsBuf bytes.Buffer
	tJS.Execute(&jsBuf, data)
	data.JS = template.HTML(jsBuf.String())

	tCSS := texttemplate.Must(texttemplate.New("css").Delims("[[", "]]").Parse(string(data.CSS)))
	var cssBuf bytes.Buffer
	tCSS.Execute(&cssBuf, data)
	data.CSS = template.HTML(cssBuf.String())

	w.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ToHTML returns the HTML for the editor component.
func (e *Editor) ToHTML() string {
	data := e.getTemplateData()

	// Use text/template for JS and CSS to avoid auto-escaping
	tJS := texttemplate.Must(texttemplate.New("js").Delims("[[", "]]").Parse(string(data.JS)))
	var jsBuf bytes.Buffer
	tJS.Execute(&jsBuf, data)
	data.JS = template.HTML(jsBuf.String())

	tCSS := texttemplate.Must(texttemplate.New("css").Delims("[[", "]]").Parse(string(data.CSS)))
	var cssBuf bytes.Buffer
	tCSS.Execute(&cssBuf, data)
	data.CSS = template.HTML(cssBuf.String())

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "editor_component.html", data); err != nil {
		return fmt.Sprintf("Error rendering editor: %v", err)
	}

	return buf.String()
}

type templateData struct {
	ID          string
	Endpoint    string
	FlowJSON    string
	Definitions string
	CSS         template.HTML
	JS          template.HTML
}

func (e *Editor) getTemplateData() templateData {
	e.mu.RLock()
	defer e.mu.RUnlock()

	defs := make([]BlockDefinition, len(e.config.Blocks))
	for i, b := range e.config.Blocks {
		defs[i] = b.Definition()
	}

	flowJSON, _ := json.Marshal(e.flow)
	defsJSON, _ := json.Marshal(defs)

	css, _ := assets.ReadFile("internal/assets/editor_component.css")
	js, _ := assets.ReadFile("internal/assets/editor_component.js")

	// Escape single quotes in JSON for embedding in JS
	flowJSONStr := strings.ReplaceAll(string(flowJSON), "'", "\\'")
	defsJSONStr := strings.ReplaceAll(string(defsJSON), "'", "\\'")

	return templateData{
		ID:          e.id,
		Endpoint:    e.config.Endpoint,
		FlowJSON:    flowJSONStr,
		Definitions: defsJSONStr,
		CSS:         template.HTML(css),
		JS:          template.HTML(js),
	}
}

func (e *Editor) serveConfig(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	defs := make([]BlockDefinition, len(e.config.Blocks))
	for i, b := range e.config.Blocks {
		defs[i] = b.Definition()
	}

	response := struct {
		Definitions []BlockDefinition `json:"definitions"`
		Flow        []Block           `json:"flow"`
	}{
		Definitions: defs,
		Flow:        e.flow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (e *Editor) handleSave(w http.ResponseWriter, r *http.Request) {
	var newFlow []Block
	if err := json.NewDecoder(r.Body).Decode(&newFlow); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	e.mu.Lock()
	e.flow = newFlow
	e.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// GetFlow returns the current flow.
func (e *Editor) GetFlow() []Block {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.flow
}

// SetFlow sets the current flow.
func (e *Editor) SetFlow(flow []Block) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.flow = flow
}
