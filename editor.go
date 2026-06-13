package stepeditor

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

// Use text/template for all templates to avoid auto-escaping of JS/CSS and delimiters conflict with Vue
var templates = texttemplate.Must(texttemplate.New("").Delims("[[", "]]").ParseFS(assets, "internal/assets/*.html"))

// Step represents a single step in the flow.
type Step struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Data     map[string]string `json:"data"`
	Branches map[string][]Step `json:"branches,omitempty"`
}

// StepDefinition provides information about a step type.
type StepDefinition struct {
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"` // Bootstrap icon class
	DefaultData map[string]string `json:"defaultData"`
	BranchNames []string          `json:"branchNames,omitempty"`
}

// CustomStep is an interface that users can implement to provide their own steps.
type CustomStep interface {
	StepDefinition() StepDefinition
}

// Config is the configuration for the Step Flow Editor.
type Config struct {
	ID              string
	Endpoint        string
	InitialFlow     []Step
	StepDefinitions []CustomStep
}

// Editor handles the HTTP requests for the step flow editor.
type Editor struct {
	id     string
	config Config
	flow   []Step
	mu     sync.RWMutex
}

// New creates a new Step Flow Editor.
func New(config Config) *Editor {
	// Ensure endpoint doesn't end with slash for consistency
	config.Endpoint = strings.TrimSuffix(config.Endpoint, "/")

	id := config.ID
	if id == "" {
		b := make([]byte, 4)
		rand.Read(b)
		id = fmt.Sprintf("editor_%x", b)
	}

	if config.InitialFlow == nil {
		config.InitialFlow = make([]Step, 0)
	}

	return &Editor{
		id:     id,
		config: config,
		flow:   config.InitialFlow,
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

	w.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ToHTML returns the HTML for the editor component.
func (e *Editor) ToHTML() string {
	data := e.getTemplateData()

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "editor_component.html", data); err != nil {
		return fmt.Sprintf("Error rendering editor: %v", err)
	}

	return buf.String()
}

type templateData struct {
	ID              string
	Endpoint        string
	FlowJSON        string
	StepDefinitions string
	CSS             template.HTML
	JS              template.HTML
}

func (e *Editor) getTemplateData() templateData {
	e.mu.RLock()
	defer e.mu.RUnlock()

	defs := make([]StepDefinition, len(e.config.StepDefinitions))
	for i, b := range e.config.StepDefinitions {
		defs[i] = b.StepDefinition()
	}

	flowJSON, _ := json.Marshal(e.flow)
	defsJSON, _ := json.Marshal(defs)

	cssRaw, _ := assets.ReadFile("internal/assets/editor_component.css")
	jsRaw, _ := assets.ReadFile("internal/assets/editor_component.js")

	interpolationData := struct {
		ID       string
		Endpoint string
	}{
		ID:       e.id,
		Endpoint: e.config.Endpoint,
	}

	// Interpolate ID and Endpoint into the JS and CSS assets before they are injected
	tJS := texttemplate.Must(texttemplate.New("js").Delims("[[", "]]").Parse(string(jsRaw)))
	var jsBuf bytes.Buffer
	tJS.Execute(&jsBuf, interpolationData)

	tCSS := texttemplate.Must(texttemplate.New("css").Delims("[[", "]]").Parse(string(cssRaw)))
	var cssBuf bytes.Buffer
	tCSS.Execute(&cssBuf, interpolationData)

	// Escape single quotes in JSON for embedding in JS strings if necessary
	flowJSONStr := strings.ReplaceAll(string(flowJSON), "'", "\\'")
	defsJSONStr := strings.ReplaceAll(string(defsJSON), "'", "\\'")

	return templateData{
		ID:              e.id,
		Endpoint:        e.config.Endpoint,
		FlowJSON:        flowJSONStr,
		StepDefinitions: defsJSONStr,
		CSS:             template.HTML(cssBuf.String()),
		JS:              template.HTML(jsBuf.String()),
	}
}

func (e *Editor) serveConfig(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	defs := make([]StepDefinition, len(e.config.StepDefinitions))
	for i, b := range e.config.StepDefinitions {
		defs[i] = b.StepDefinition()
	}

	response := struct {
		StepDefinitions []StepDefinition `json:"stepDefinitions"`
		Flow            []Step           `json:"flow"`
	}{
		StepDefinitions: defs,
		Flow:            e.flow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (e *Editor) handleSave(w http.ResponseWriter, r *http.Request) {
	var newFlow []Step
	if err := json.NewDecoder(r.Body).Decode(&newFlow); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	e.mu.Lock()
	e.flow = newFlow
	e.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetFlow returns the current flow.
func (e *Editor) GetFlow() []Step {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.flow
}

// SetFlow sets the current flow.
func (e *Editor) SetFlow(flow []Step) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.flow = flow
}
