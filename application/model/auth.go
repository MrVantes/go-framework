package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
)

type UserJwtCustomClaims struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	UserID           int32  `json:"user_id"`
	OrganizationID   int32  `json:"org_id"`
	OrganizationName string `json:"org_name"`
	jwt.StandardClaims
}

type SigninRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type SigninResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
}

type AppUser struct {
	UserID      int32  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func NewAppUser(p *db.AppUser) *AppUser {
	return &AppUser{
		UserID:      p.UserID,
		Username:    p.Username,
		DisplayName: p.DisplayName,
		Email:       p.Email,
	}
}

type LoginHistoryParams struct {
	UserID      int32     `json:"user_id"`
	Username    string    `json:"username"`
	LoginTime   time.Time `json:"login_time"`
	LoginStatus bool      `json:"login_status"`
}

func NewCreateLoginHistoryParams(p *LoginHistoryParams) *db.CreateLoginHistoryParams {
	return &db.CreateLoginHistoryParams{
		UserID:      pgtype.Int4{Int32: p.UserID, Valid: true},
		Username:    pgtype.Text{String: p.Username, Valid: true},
		LoginTime:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		LoginStatus: pgtype.Bool{Bool: p.LoginStatus, Valid: true},
	}
}
