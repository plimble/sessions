package mongo

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func tearDown(store *Store) {
	store.getC().DropCollection()
}

func TestNewNotFoundName(t *testing.T) {
	sess, err := mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		panic("MONGO CANT CONNECT")
	}
	defer sess.Close()

	s := NewStore(sess, "/", 38600)
	r, _ := http.NewRequest("GET", "/", nil)

	// ================ NEW ================
	session, err := s.New(r, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", session.Name())
	assert.True(t, session.IsNew)
	assert.NotNil(t, session.ID)
	// =====================================
	tearDown(s)
}

func TestNewFoundName(t *testing.T) {
	sess, err := mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		panic("MONGO CANT CONNECT")
	}
	defer sess.Close()

	s := NewStore(sess, "/", 38600)
	r, _ := http.NewRequest("GET", "/", nil)
	c := http.Cookie{Name: "test"}
	r.AddCookie(&c)

	// ================ NEW ================
	session, err := s.New(r, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", session.Name())
	assert.True(t, session.IsNew)
	assert.NotNil(t, session.ID)
	// =====================================

	// ================ SAVE ===============
	w := httptest.NewRecorder()
	err = s.Save(r, w, session)
	assert.NoError(t, err)
	assert.NotEqual(t, "test=", w.HeaderMap["Set-Cookie"])
	// =====================================

	// ============== NEW AGAIN ============

	r, _ = http.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{
		Name:     "test",
		Value:    session.ID,
		Path:     s.options.Path,
		Domain:   s.options.Domain,
		MaxAge:   s.options.MaxAge,
		Secure:   s.options.Secure,
		HttpOnly: s.options.HttpOnly,
	}
	r.AddCookie(cookie)

	session, err = s.New(r, "test")

	assert.NoError(t, err)
	assert.Equal(t, "test", session.Name())
	assert.False(t, session.IsNew)
	// =====================================
}
