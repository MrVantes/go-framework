package security

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/arifnurdiansyah92/go-boilerplate/application/config"
	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
	"github.com/arifnurdiansyah92/go-boilerplate/application/pkg/registry"
	"github.com/arifnurdiansyah92/go-boilerplate/application/pkg/response"
)

type UserApp struct {
	q   *db.Queries
	cfg *config.Config
}

func NewUserApp(dbPool *pgxpool.Pool, cfg *config.Config) *UserApp {
	return &UserApp{
		q:   db.New(dbPool),
		cfg: cfg,
	}
}

func (a *UserApp) RegisterRoutes(e *echo.Group, m *registry.Module) {
	// Define actions with handlers
	actions := []registry.Action{
		{Name: "CreateUser", Method: "POST", Path: "/users/create", Handler: a.CreateUser},
		{Name: "GetUser", Method: "GET", Path: "/users/:id", Handler: a.GetUser},
		{Name: "ListUsers", Method: "GET", Path: "/users", Handler: a.ListUsers},
		{Name: "UpdateUser", Method: "PUT", Path: "/users/:id", Handler: a.UpdateUser},
		{Name: "DeleteUser", Method: "DELETE", Path: "/users/:id", Handler: a.DeleteUser},
		{Name: "Profile", Method: "GET", Path: "/me", Handler: a.GetProfile},
	}

	// Add app to module
	m.Apps = append(m.Apps, registry.App{
		Name:    "Users",
		Actions: actions,
	})

	// Dynamically register routes
	registry.RegisterRoutes(e, actions)
}

func (a *UserApp) CreateUser(c echo.Context) error {
	return nil
}

func (a *UserApp) GetUser(c echo.Context) error {
	return nil
}

func (a *UserApp) ListUsers(c echo.Context) error {
	return nil
}

func (a *UserApp) UpdateUser(c echo.Context) error {
	return nil
}

func (a *UserApp) DeleteUser(c echo.Context) error {
	return nil
}

func (a *UserApp) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(int32)
	user, err := a.q.GetUserByID(context.Background(), userID)
	if err != nil {
		return response.Error(c, http.StatusUnprocessableEntity, err)
	}

	user.Password = ""

	return c.JSON(http.StatusOK, echo.Map{
		"data": user,
	})
}
