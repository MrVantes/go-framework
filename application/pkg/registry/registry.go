package registry

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Action struct {
	Name    string           `json:"name"`
	Method  string           `json:"method"`
	Path    string           `json:"path"`
	Handler echo.HandlerFunc `json:"-"` // Don't serialize the handler
}

type App struct {
	Name    string   `json:"name"`
	Actions []Action `json:"actions"`
}

type Module struct {
	Name string `json:"name"`
	Apps []App  `json:"apps"`
}

type Registry struct {
	Modules []Module `json:"modules"`
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) AddModule(m Module) {
	r.Modules = append(r.Modules, m)
}

func (r *Registry) GetModule(name string) *Module {
	for _, m := range r.Modules {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

func (r *Registry) GetModules() []Module {
	return r.Modules
}

func (r *Registry) AddApp(moduleName string, a App) {
	for i, m := range r.Modules {
		if m.Name == moduleName {
			r.Modules[i].Apps = append(r.Modules[i].Apps, a)
		}
	}
}

func (r *Registry) GetApp(moduleName, appName string) *App {
	for _, m := range r.Modules {
		if m.Name == moduleName {
			for _, a := range m.Apps {
				if a.Name == appName {
					return &a
				}
			}
		}
	}
	return nil
}

func (r *Registry) AddAction(moduleName, appName string, a Action) {
	for i, m := range r.Modules {
		if m.Name == moduleName {
			for j, app := range m.Apps {
				if app.Name == appName {
					r.Modules[i].Apps[j].Actions = append(r.Modules[i].Apps[j].Actions, a)
				}
			}
		}
	}
}

func (a *Registry) GetModulesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, a.GetModules())
}

// RegisterRoutes dynamically registers routes for a given Echo group based on a slice of Actions
func RegisterRoutes(e *echo.Group, actions []Action) {
	for _, action := range actions {
		switch action.Method {
		case "POST":
			e.POST(action.Path, action.Handler)
		case "GET":
			e.GET(action.Path, action.Handler)
		case "PUT":
			e.PUT(action.Path, action.Handler)
		case "DELETE":
			e.DELETE(action.Path, action.Handler)
			// Add more cases if needed
		}
	}
}
