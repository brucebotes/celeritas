package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieLifeTime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
}

func (c *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// how long should persist last?
	minutes, err := strconv.Atoi(c.CookieLifeTime)
	if err != nil {
		minutes = 60
	}

	// should cookies persist
	if strings.ToLower(c.CookiePersist) == "true" {
		persist = true
	}

	//must cookies be secure
	if strings.ToLower(c.CookieSecure) == "true" {
		secure = true
	}

	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = c.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = c.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	switch strings.ToLower(c.SessionType) {
	case "redis":
		session.Store = redisstore.New(c.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(c.DBPool)
	case "postgres", "postgresql":

		session.Store = postgresstore.New(c.DBPool)
	case "sqlite", "sqlite3":
		/*
		* --------------------------------------------------
		* TODO: Please update the sqlite3 database (./system)
		*		to refect the following
		* there has been a slight change in the newer version
		* requirements for sqlite3store.
		* --------------------------------------------------
		*    CREATE TABLE sessions (
		*      token TEXT PRIMARY KEY,
		*      data BLOB NOT NULL,
		*      expiry REAL NOT NULL
		*    );
		*
		*    CREATE INDEX sessions_expiry_idx ON sessions(expiry);
		 */
		session.Store = sqlite3store.New(c.DBPool)
	default:
		// cookie
	}
	return session
}
