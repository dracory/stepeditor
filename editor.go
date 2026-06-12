package stepfloweditor

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

//go:embed internal/assets/index.html
var assets embed.FS

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
	Endpoint string
	Blocks   []CustomBlock
}

// Editor handles the HTTP requests for the step flow editor.
type Editor struct {
	config NewConfig
	flow   []Block
	mu     sync.RWMutex
}

// New creates a new Step Flow Editor.
func New(config NewConfig) *Editor {
	// Ensure endpoint doesn't end with slash for consistency
	config.Endpoint = strings.TrimSuffix(config.Endpoint, "/")
	return &Editor{
		config: config,
		flow:   make([]Block, 0),
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
	index, err := assets.ReadFile("internal/assets/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(index)
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
