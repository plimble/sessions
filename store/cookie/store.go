package cookie

import (
	"bytes"
	"encoding/base64"
	"github.com/plimble/sessions"
	"net/http"
)

type CookieStore struct{}

func NewCookieStore() *CookieStore {
	return &CookieStore{}
}

func (s *CookieStore) Get(id string, buf *bytes.Buffer) error {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return err
	}

	if _, err := buf.Write(b); err != nil {
		return err
	}

	return nil
}

func (s *CookieStore) Delete(id string) error {
	return nil
}

func (s *CookieStore) Save(session *sessions.Session, buf *bytes.Buffer, w http.ResponseWriter) error {
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}
