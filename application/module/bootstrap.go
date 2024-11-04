package module

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/arifnurdiansyah92/go-boilerplate/application/config"
	"github.com/arifnurdiansyah92/go-boilerplate/application/db"
)

type OrgBootstrap struct {
	q   *db.Queries
	cfg *config.Config
}

func NewOrgBootstrap(dbPool *pgxpool.Pool, cfg *config.Config) *OrgBootstrap {
	q := db.New(dbPool)
	return &OrgBootstrap{
		q:   q,
		cfg: cfg,
	}
}

func (h *OrgBootstrap) InitialOrg() {
	orgData := h.cfg.Bootstrap.Initial

	// Do not allow empty initial organizations
	if orgData.OrgName == "" || orgData.AdminUsername == "" || orgData.AdminPassword == "" || orgData.AdminEmail == "" || orgData.AdminDisplayname == "" {
		log.Fatal().Msg("Initial organizations data is empty")
		return
	}

	ctx := context.Background()
	_, err := h.q.GetOrganizationByName(ctx, orgData.OrgName)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {

			createOrganizationParams := &db.CreateOrganizationParams{
				OrganizationName:        orgData.OrgName,
				OrganizationDescription: orgData.OrgDescription,
			}
			org, err := h.q.CreateOrganization(ctx, createOrganizationParams)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create initial organizations")
			}
			log.Info().Msg("Initial organizations created")

			// Create initial organizations admin user
			password := orgData.AdminPassword
			hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to generate password hash")
			}

			createUserParams := &db.CreateUserParams{
				Username:       orgData.AdminUsername,
				Password:       string(hashedBytes),
				Email:          orgData.AdminEmail,
				DisplayName:    orgData.AdminDisplayname,
				OrganizationID: org.OrganizationID,
			}

			_, err = h.q.CreateUser(ctx, createUserParams)
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to create initial org admin user")
			}
			log.Info().Msg("Initial org admin user created")

		} else {
			log.Fatal().Err(err).Msg("Unable to get initial organizations")
		}
	} else {
		log.Info().Msg("Initial organizations found")
	}
}
