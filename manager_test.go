package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SessionManagerSuite struct {
	suite.Suite
	m *SessionManager
}

func TestSessionManagerSuite(t *testing.T) {
	suite.Run(t, &SessionManagerSuite{})
}

func (t *SessionManagerSuite) SetupSuite() {
	store := NewMemoryStore()
	t.m = New(10000, store, nil)
}

func (t *SessionManagerSuite) TestSession() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	sessions := t.m.GetSessions(req)
	session, err := sessions.Get("test")
	t.Error(err)

	session.Set("foo", "abc")
	session.Set("baz", 123)
	session.Set("test", []string{"1", "2", "3"})
	session.Set("test2", []int{1, 2, 3})
	session.Set("bool", true)
	session.Set("float", 10.11)
	session.Set("floats", []float64{10.11, 12.13})

	err = sessions.Save(w)
	t.NoError(err)
	t.True(session.IsNew)

	t.m.Close(sessions)

	hdr := w.Header()
	cookies, ok := hdr["Set-Cookie"]
	t.True(ok)
	t.Len(cookies, 1)

	req, _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Add("Cookie", cookies[0])
	w = httptest.NewRecorder()

	sessions = t.m.GetSessions(req)
	session, err = sessions.Get("test")
	t.NoError(err)

	t.False(session.IsNew)
	t.Equal("abc", session.GetString("foo", ""))
	t.Equal(123, session.GetInt("baz", 0))
	t.Equal([]string{"1", "2", "3"}, session.GetStrings("test", nil))
	t.Equal([]int64{1, 2, 3}, session.GetInts("test2", nil))
	t.Equal(true, session.GetBool("bool", false))
	t.Equal(10.11, session.GetFloat("float", 0))
	t.Equal([]float64{10.11, 12.13}, session.GetFloats("floats", nil))
}

func (t *SessionManagerSuite) TestFlashes() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	sessions := t.m.GetSessions(req)
	session, err := sessions.Get("test")
	t.Error(err)

	session.AddFlash("1")
	session.AddFlash("2")
	session.AddFlash("3")

	err = sessions.Save(w)
	t.NoError(err)
	t.True(session.IsNew)

	t.m.Close(sessions)

	hdr := w.Header()
	cookies, ok := hdr["Set-Cookie"]
	t.True(ok)
	t.Len(cookies, 1)

	req, _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Add("Cookie", cookies[0])
	w = httptest.NewRecorder()

	sessions = t.m.GetSessions(req)
	session, err = sessions.Get("test")
	t.NoError(err)

	flashes := session.Flashes()

	t.False(session.IsNew)
	t.Len(flashes, 3)
	t.Equal(flashes[0], "1")
	t.Equal(flashes[1], "2")
	t.Equal(flashes[2], "3")
	_, ok = session.Values[flashesKey]
	t.False(ok)
}

func BenchmarkC2(b *testing.B) {
	store := sessions.NewCookieStore([]byte("secret-key"))

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		session, _ := store.Get(req, "session-key")
		session.Values["foo"] = "abc"
		session.Values["baz"] = 123

		session.Save(req, w)
	}
}

func BenchmarkC1(b *testing.B) {
	store := NewMemoryStore()
	m := New(10000, store, nil)

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()

		sessions := m.GetSessions(req)
		session, _ := sessions.Get("test")

		session.Set("foo", "abc")
		session.Set("baz", 123)

		sessions.Save(w)
		m.Close(sessions)
	}
}
