package mongo

import (
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
	"labix.org/v2/mgo"
	"net/http"
)

type Store struct {
	db      string
	c       string
	session *mgo.Session
	options *sessions.Options
}

func NewStore(session *mgo.Session, path string, maxAge int) *Store {
	if path == "" {
		path = "/"
	}

	if maxAge == 0 {
		maxAge = 38400
	}

	return &Store{
		db:      "store",
		c:       "sessions",
		session: session,
		options: &sessions.Options{
			Path:   path,
			MaxAge: maxAge,
		},
	}
}

func (s *Store) getC() *mgo.Collection {
	return s.session.DB(s.db).C(s.c)
}

func generateUUID() string {
	v1 := uuid.NewV1()
	return hex.EncodeToString(v1.Bytes())
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = s.options
	session.IsNew = true
	session.ID = generateUUID()

	var err error
	if c, err := r.Cookie(name); err == nil {
		if exist, err := s.getC().FindId(c.Value).Count(); err == nil {
			if exist > 0 {
				session.IsNew = false
			}
		}
	}

	return session, err
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" {
		session.ID = generateUUID()
	}

	if _, err := s.getC().UpsertId(session.ID, session); err != nil {

		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))
	return nil
}
