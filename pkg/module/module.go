package module

import "fmt"

// Module represents a Judozi module
type Module interface {
	// Name returns the module name
	Name() string
	
	// Description returns a short description
	Description() string
	
	// Category returns the module category (e.g., "privesc", "recon", "exploit")
	Category() string
	
	// Run executes the module with given arguments
	Run(args []string) error
}

// Registry holds all registered modules
type Registry struct {
	modules map[string]Module
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry
func (r *Registry) Register(m Module) {
	r.modules[m.Name()] = m
}

// Get retrieves a module by name
func (r *Registry) Get(name string) (Module, error) {
	m, ok := r.modules[name]
	if !ok {
		return nil, fmt.Errorf("module '%s' not found", name)
	}
	return m, nil
}

// List returns all registered modules
func (r *Registry) List() []Module {
	modules := make([]Module, 0, len(r.modules))
	for _, m := range r.modules {
		modules = append(modules, m)
	}
	return modules
}

// ListByCategory returns modules filtered by category
func (r *Registry) ListByCategory(category string) []Module {
	modules := make([]Module, 0)
	for _, m := range r.modules {
		if m.Category() == category {
			modules = append(modules, m)
		}
	}
	return modules
}
