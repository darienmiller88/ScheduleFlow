package sessionmanager

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
)

func NewSessionManager(db *sqlx.DB) *scs.SessionManager{
	sessionManager := scs.New()

	sessionManager.Lifetime = 168 * time.Hour
	sessionManager.Store = postgresstore.New(db.DB)
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true

	return sessionManager
}