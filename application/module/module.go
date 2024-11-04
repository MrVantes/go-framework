package module

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/arifnurdiansyah92/go-boilerplate/application/config"
	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
	internalMiddleware "github.com/arifnurdiansyah92/go-boilerplate/application/middleware"
	"github.com/arifnurdiansyah92/go-boilerplate/application/module/security"
	"github.com/arifnurdiansyah92/go-boilerplate/application/pkg/registry"
)

func InitApps(e *echo.Echo, dbPool *pgxpool.Pool, cfg *config.Config) {

	// auth routes is for signin, signout
	a := e.Group("/auth")
	authHandler := NewAuthHandler(dbPool, cfg)
	authHandler.SetRoutes(a)

	// Registry for all apps
	r := registry.NewRegistry()

	// user routes will check for bearer token and authorization
	u := e.Group("")

	u.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.UserJWT.SigningKey),
		TokenLookup: "header:Authorization,query:access_token",
	}))

	authzMiddleware := &internalMiddleware.Authorization{
		Q: db.New(dbPool),
	}
	u.Use(authzMiddleware.Check)

	// // ====================
	// assets module
	// assetModule := &registry.Module{
	// 	Name: "Assets",
	// }

	// assets app
	// assetApp := asset.NewAssetApp(dbPool, cfg)
	// assetApp.RegisterRoutes(u, assetModule)

	// other assets apps here ...

	// register module
	// r.AddModule(*assetModule)

	// ====================
	// user module
	securityModule := &registry.Module{
		Name: "Security",
	}

	// user app
	userApp := security.NewUserApp(dbPool, cfg)
	userApp.RegisterRoutes(u, securityModule)

	// other security apps here ...

	// register module
	r.AddModule(*securityModule)

	// ====================
	// register registry
	u.GET("/modules", r.GetModulesHandler)

	// ====================
	// bootstrap initial org
	orgBootstrap := NewOrgBootstrap(dbPool, cfg)
	orgBootstrap.InitialOrg()
}
