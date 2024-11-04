package module

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/arifnurdiansyah92/go-boilerplate/application/config"
	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
	"github.com/arifnurdiansyah92/go-boilerplate/application/model"
	"github.com/arifnurdiansyah92/go-boilerplate/application/pkg/response"
)

type AuthHandler struct {
	q   *db.Queries
	cfg *config.Config
}

func NewAuthHandler(dbPool *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	h := &AuthHandler{
		q:   db.New(dbPool),
		cfg: cfg,
	}

	return h
}

// SetRoutes ...
func (h *AuthHandler) SetRoutes(e *echo.Group) {
	e.POST("/signin", h.signin)
	e.POST("/signout", h.signout)
}

func (h *AuthHandler) signin(c echo.Context) error {
	requestData := &model.SigninRequest{}
	err := c.Bind(requestData)
	if err != nil {
		return response.Error(c, http.StatusUnprocessableEntity, err)
	}

	//init login history params
	loginHistory := &model.LoginHistoryParams{}

	ctx := context.Background()
	user, err := h.q.GetUserByName(ctx, requestData.Username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.Error(c, http.StatusUnauthorized, errors.New("signin: invalid username or password"))
		}
		return response.Error(c, http.StatusInternalServerError, err)
	}

	if !user.IsActive.Bool {
		return response.Error(c, http.StatusUnauthorized, errors.New("signin: account inactive"))
	}

	lastLoginAttempt, err := h.q.GetLastLogin(ctx, pgtype.Text{String: user.Username, Valid: true})
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			return response.Error(c, http.StatusInternalServerError, err)
		}
	}

	if user.IsLocked.Bool && time.Since(lastLoginAttempt.Time) < time.Hour {
		return response.Error(c, http.StatusUnauthorized, errors.New("signin: account locked, try again later"))
	}

	// Default value for login history
	loginHistory.UserID = user.UserID
	loginHistory.Username = user.Username
	loginHistory.LoginStatus = true

	org, err := h.q.GetOrganization(ctx, user.OrganizationID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))
	if err != nil {
		// Set Login Status False
		loginHistory.LoginStatus = false
		errLog := h.q.CreateLoginHistory(ctx, model.NewCreateLoginHistoryParams(loginHistory))
		if errLog != nil {
			return response.Error(c, http.StatusInternalServerError, errLog)
		}
		return response.Error(c, http.StatusUnauthorized, errors.New("signin: invalid username or password"))
	}

	// Set JWT claims
	duration, err := time.ParseDuration(h.cfg.UserJWT.Duration)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}
	claims := &model.UserJwtCustomClaims{
		user.Username,
		user.Email,
		user.UserID,
		org.OrganizationID,
		org.OrganizationName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.cfg.UserJWT.SigningKey))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}

	errLog := h.q.CreateLoginHistory(ctx, model.NewCreateLoginHistoryParams(loginHistory))
	if errLog != nil {
		return response.Error(c, http.StatusInternalServerError, errLog)
	}

	return c.JSON(http.StatusOK, model.SigninResponse{
		Token: t,
	})
}

func (h *AuthHandler) signout(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
