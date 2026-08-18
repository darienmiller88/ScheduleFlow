package sessionmanager

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/alexedwards/scs/postgresstore"
)

func NewSessionManager(db *sqlx.DB) *scs.SessionManager{
	sessionManager := scs.New()

	sessionManager.Lifetime = 168 * time.Hour
	sessionManager.Store = postgresstore.New(db.DB)

	return sessionManager
}