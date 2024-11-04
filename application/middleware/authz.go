package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
	"github.com/arifnurdiansyah92/go-boilerplate/application/pkg/response"
)

type Authorization struct {
	Q *db.Queries
}

func (m *Authorization) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(float64)
		username := claims["username"].(string)
		orgID := claims["org_id"].(float64)
		orgName := claims["org_name"].(string)

		// Check that user exists
		ctx := context.Background()
		getUserParams := &db.GetUserByOrgParams{
			UserID:         int32(userID),
			OrganizationID: int32(orgID),
		}
		_, err := m.Q.GetUserByOrg(ctx, getUserParams)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return response.Error(c, http.StatusUnauthorized, errors.New("authz: invalid username"))
			}
			return response.Error(c, http.StatusInternalServerError, err)
		}

		// Check user access to resource
		ok := false

		// TODO: Implement access control logic here
		// For example, check if user has access to a specific resource
		// by checking if the user is in a specific security group
		// for now we'll just set ok to true
		ok = true

		// If user does not have access to resource, return unauthorized
		if !ok {
			return response.Error(c, http.StatusUnauthorized, errors.New("authz: unauthorized access to resource"))
		}

		// Set user_id, username, org_id, and org_name in context for use in route handlers
		c.Set("user_id", int32(userID))
		c.Set("username", username)
		c.Set("org_id", int32(orgID))
		c.Set("org_name", orgName)

		return next(c)
	}
}
